package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/lib/pq"
	"net/http"
	"time"
)

// CreateUserRequest represents a request to create a new user.
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=40"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6,max=40"`
	Email     string `json:"email" binding:"required,email"`
}

// CreateUserResponse represents a response from a create user request.
type CreateUserResponse struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// NewUserResponse creates a new CreateUserResponse from a db.User.
func NewUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		ID:                user.ID,
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// CreateNewUser creates a new user account.
func (s *Server) CreateNewUser(context *gin.Context) {
	var req CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	createUserArg := db.CreateUserParams{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	createRoleArg := db.CreateRoleParams{
		Name: "user",
	}

	user, err := s.store.CreateUserTx(context, db.CreateUserTxParams{
		CreateUserParams: createUserArg,
		CreateRoleParams: createRoleArg,
	})
	if err != nil {
		if pgErr, err := err.(*pq.Error); err {
			switch pgErr.Code.Name() {
			case "unique_violation":
				context.JSON(400, helpers.ErrorResponse(pgErr))
				return
			}
		}
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	rsp := NewUserResponse(user.User)
	context.JSON(http.StatusOK, rsp)
}

// ListUsersRequest represents a request to list users.
type ListUsersRequest struct {
	PageSize int32 `form:"page_size" binding:"required,min=1,max=25"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
}

// UserResponse represents a response from a list users request.
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ListUsers lists all users.
func (s *Server) ListUsers(context *gin.Context) {
	var req ListUsersRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := s.store.ListUsers(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	rsp := make([]UserResponse, len(users))
	for i, user := range users {
		rsp[i] = UserResponse{
			ID:        user.ID.String(),
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	context.JSON(200, rsp)
}

// GetUserRequest represents a request to get a user.
type GetUserRequest struct {
	ID string `uri:"id" binding:"required,uuid4"`
}

// GetUserResponse represents a response from a get user request.
type GetUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUser gets a user.
func (s *Server) GetUser(context *gin.Context) {
	var req GetUserRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	userID := uuid.MustParse(req.ID)

	user, err := s.store.GetUser(context, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	rsp := GetUserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	context.JSON(http.StatusOK, rsp)
}
