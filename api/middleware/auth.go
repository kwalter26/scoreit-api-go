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
	AuthorizationTypeBearer = "bearer"
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

//func Auth0Middleware(config util.Config) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		issuerURL, err := url.Parse("https://" + config.Auth0Domain + "/")
//		if err != nil {
//			log.Fatal().Err(err).Msg("failed to parse the issuer url")
//		}
//
//		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
//
//		jwtValidator, _ := validator.New(
//			provider.KeyFunc,
//			validator.RS256,
//			issuerURL.String(),
//			[]string{config.Auth0Audience},
//		)
//		jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
//
//		//err := jwtMiddleware.CheckJWT(c.Writer, c.Request)
//		if err != nil {
//			// Handle authentication error
//			c.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		// Continue to the next middleware/handler
//		c.Next()
//	}
//}

func GetAuthorizationPayload(c *gin.Context) *token.Payload {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return nil
	}
	return payload.(*token.Payload)
}
