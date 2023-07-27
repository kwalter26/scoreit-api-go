package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"net/http"
)

// PasetoAuthorizer stores the casbin handler
type PasetoAuthorizer struct {
	enforcer *casbin.Enforcer
}

// NewAuthorizeMiddleware returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizeMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	a := &PasetoAuthorizer{enforcer: e}

	return func(c *gin.Context) {
		if !a.CheckPermission(a.GetUser(c), c.Request) {
			a.RequirePermission(c)
		}
	}
}

func (a *PasetoAuthorizer) GetUser(c *gin.Context) string {
	ap, exists := c.Get(AuthorizationPayloadKey)
	var user string
	if !exists {
		return ""
	} else {
		p := ap.(*token.Payload)
		user = p.UserID.String()
		var groups [][]string
		for _, g := range p.Permissions {
			groups = append(groups, []string{user, string(g)})
		}
		_, _ = a.enforcer.AddGroupingPolicies(groups)
		return user
	}
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *PasetoAuthorizer) CheckPermission(user string, r *http.Request) bool {
	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		panic(err)
	}

	return allowed
}

// RequirePermission returns the 403 Forbidden to the client
func (a *PasetoAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
