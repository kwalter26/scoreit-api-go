package middleware

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Auth0Authorizer struct declaration with a description of its role.
type Auth0Authorizer struct {
	enforcer *casbin.Enforcer // A field of type '*casbin.Enforcer'
}

// NewAuthorizeMiddleware function declaration.
// A description of its behavior is provided.
// This function returns a Gin middleware function.
//
// Example usage and a walkthrough is also included for a more in-depth understanding.
func NewAuthorizeMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	a := &Auth0Authorizer{enforcer: e}

	return func(c *gin.Context) {
		user, groups := a.GetUser(c)
		if !a.CheckPermission(user, groups, c.Request) {
			a.RequirePermission(c)
		}
	}
}

// GetUser function declaration.
// Explains the retrieval of user and associated groups from gin context.
// Also describes possible return scenarios from the function.
func (a *Auth0Authorizer) GetUser(c *gin.Context) (string, []security.Role) {
	key, exists := c.Get(AuthorizationPayloadKey)
	user := ""
	var groups []security.Role
	if !exists {
		return user, groups
	} else {
		validatedClaims := key.(*validator.ValidatedClaims)
		user = validatedClaims.RegisteredClaims.Subject
		if validatedClaims.CustomClaims == nil {
			return user, groups
		}
		customClaims := validatedClaims.CustomClaims.(*CustomClaims)
		for _, group := range customClaims.Roles {
			groups = append(groups, group)
		}
		return user, groups
	}
}

// CheckPermission method declaration.
// This method checks if a user or group has necessary permissions.
func (a *Auth0Authorizer) CheckPermission(user string, groups []security.Role, r *http.Request) bool {
	method := r.Method
	path := r.URL.Path

	if allowed, err := a.enforcer.Enforce(user, path, method); err != nil {
		return false
	} else if allowed {
		log.Debug().Str("user", user).Str("method", method).Str("path", path).Msgf("user allowed to access resource")
		return true
	}

	for _, g := range groups {
		if allowed, err := a.enforcer.Enforce(string(g), path, method); err != nil {
			log.Error().Err(err).Str("user", user).Msg("error checking permissions of user")
			return false
		} else if allowed {
			log.Debug().Str("group", user).Str("method", method).Str("path", path).Msgf("group allowed to access resource")
			return true
		}
	}
	return false
}

// RequirePermission method declaration.
// This method aborts the current HTTP request and sends a 403 Forbidden status to the client.
func (a *Auth0Authorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
