// Package middleware provides middlewares for various functionalities like JWT validation, etc.
package middleware

import (
	"context"
	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"time"
)

// CustomClaims contains custom data we want from the token.
// It stores Scope and Roles, and satisfies validator.CustomClaims interface.
type CustomClaims struct {
	Scope string          `json:"scope"`
	Roles []security.Role `json:"scoreit.us.auth0.com/roles"`
}

// NewAuth0JwtValidator creates a new JWT Validator configured to use Auth0's settings based on the provided configuration.
// It fetches keys from JWKS provider and validates the JWT signatures.
func NewAuth0JwtValidator(config util.Config) *validator.Validator {
	issuerURL, err := url.Parse("https://" + config.Auth0Domain + "/")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse the issuer url")
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Auth0Audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
		validator.WithAllowedClockSkew(time.Minute),
	)
	// check error
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set up the jwt validator")
	}
	return jwtValidator
}

// Validate is a dummy method which satisfies the validator.CustomClaims interface.
// However, it doesn't perform any validations.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// CheckJWT is a middleware that checks and validates the incoming JWT in the HTTP request header.
// In case of any errors in validation, it aborts the request and sends an 'Unauthorized' status.
func CheckJWT(v *validator.Validator) gin.HandlerFunc {

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Msg("Encountered error while validating JWT")
	}

	middleware := jwtmiddleware.New(
		v.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			token := ctx.Request.Context().Value(jwtmiddleware.ContextKey{})
			ctx.Set(AuthorizationPayloadKey, token)
			ctx.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}
