package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

// LoggerMiddleware is a middleware function that logs information about incoming requests.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer
		c.Next()

		logParams := generateLogParams(c, start)
		logEvent := createLogEvent(c, logParams)

		logEvent.Msg(logParams.ErrorMessage)
	}
}

// generateLogParams is a function that generates log parameters based on the given gin.Context and start time.
// It returns a gin.LogFormatterParams struct containing the following information:
// - TimeStamp: current time when the log parameters are generated.
// - Latency: time elapsed since the start time.
// - ClientIP: IP address of the client.
// - Method: HTTP request method.
// - StatusCode: HTTP status code of the response.
// - ErrorMessage: string representation of any errors that occurred during the request.
// - BodySize: size of the response body.
// - Path: URL path of the request.
//
// If the latency exceeds one minute, it will be truncated to the nearest second.
//
// Example usage:
//
//	start := time.Now()
//	c.Next()
//	logParams := generateLogParams(c, start)
//	logEvent := createLogEvent(c, logParams)
//	logEvent.Msg(logParams.ErrorMessage)
func generateLogParams(c *gin.Context, start time.Time) gin.LogFormatterParams {
	path := c.Request.URL.Path

	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	logParams := gin.LogFormatterParams{
		TimeStamp:    time.Now(),
		Latency:      time.Since(start),
		ClientIP:     c.ClientIP(),
		Method:       c.Request.Method,
		StatusCode:   c.Writer.Status(),
		ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
		BodySize:     c.Writer.Size(),
		Path:         path,
	}

	if logParams.Latency > time.Minute {
		logParams.Latency = logParams.Latency.Truncate(time.Second)
	}

	return logParams
}

// createLogEvent is a function that creates a log event based on the given context and log parameters.
func createLogEvent(c *gin.Context, logParams gin.LogFormatterParams) *zerolog.Event {
	if c.Writer.Status() >= 500 {
		return log.Error().
			Str("client_ip", logParams.ClientIP).
			Str("method", logParams.Method).
			Int("status_code", logParams.StatusCode).
			Int("body_size", logParams.BodySize).
			Str("path", logParams.Path).
			Str("latency", logParams.Latency.String())
	}

	return log.Info().
		Str("client_ip", logParams.ClientIP).
		Str("method", logParams.Method).
		Int("status_code", logParams.StatusCode).
		Int("body_size", logParams.BodySize).
		Str("path", logParams.Path).
		Str("latency", logParams.Latency.String())
}
