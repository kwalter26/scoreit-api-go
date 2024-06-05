package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
			return
		}

		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}

func GetAuthorizationPayload(c *gin.Context) *token.Payload {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return nil
	}
	return payload.(*token.Payload)
}
