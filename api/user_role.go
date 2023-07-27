package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"net/http"
)

type GetUserRolesRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetUserRolesResponse struct {
	Name string `json:"name"`
}

func (s *Server) GetUserRoles(context *gin.Context) {
	var req GetUserRolesRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	userId := uuid.MustParse(req.Id)

	roles, err := s.store.GetRoles(context, userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	var resp []GetUserRolesResponse
	for _, role := range roles {
		resp = append(resp, GetUserRolesResponse{Name: role.Name})
	}

	if resp == nil {
		resp = []GetUserRolesResponse{}
	}

	context.JSON(http.StatusOK, resp)
}

type CreateUserRolesRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CreateUserRolesRequestBody struct {
	Name string `json:"name" binding:"required"`
}

type CreateUserRolesResponse struct {
	Name string `json:"name"`
}

func (s *Server) CreateUserRole(context *gin.Context) {
	var req CreateUserRolesRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	userId := uuid.MustParse(req.Id)

	var body CreateUserRolesRequestBody
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	arg := db.CreateRoleParams{
		Name:   body.Name,
		UserID: userId,
	}

	role, err := s.store.CreateRole(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	resp := CreateUserRolesResponse{Name: role.Name}

	context.JSON(http.StatusOK, resp)
}

type ListUserRolesRequest struct {
	PageSize int32 `form:"page_size,default=5" binding:"number,max=100,min=1"`
	PageID   int32 `form:"page_id,default=1" binding:"number,min=1"`
}

func (s *Server) ListUserRoles(context *gin.Context) {
	var req ListUserRolesRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	arg := db.ListRolesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := s.store.ListRoles(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.JSON(http.StatusOK, roles)

}
