package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggerMiddleware(t *testing.T) {
	// Create new test router
	router := gin.New()

	router.Use(LoggerMiddleware())
	router.GET("/example", func(c *gin.Context) {
		c.String(http.StatusOK, "example")
	})

	// table driven tests
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "normalRequest",
			path:     "/example",
			expected: http.StatusOK,
		},
		{
			name:     "nonExistent",
			path:     "/nonExistentPath",
			expected: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request to pass to our router
			req, _ := http.NewRequest("GET", tc.path, nil)

			// Test router serves the request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expected {
				t.Errorf("expected status: %v, got: %v", tc.expected, w.Code)
			}
		})
	}
}

func TestCreateLogEvent(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		expErr    bool
		logParams gin.LogFormatterParams
	}{
		{
			name:   "BadRequest",
			status: 500,
			expErr: true,
			logParams: gin.LogFormatterParams{
				ClientIP:   "192.0.0.1",
				Method:     "GET",
				StatusCode: 500,
				BodySize:   15,
				Path:       "/api/test",
				Latency:    time.Second,
			},
		},
		{
			name:   "SuccessfulReq",
			status: 200,
			expErr: false,
			logParams: gin.LogFormatterParams{
				ClientIP:   "192.0.0.1",
				Method:     "GET",
				StatusCode: 200,
				BodySize:   15,
				Path:       "/api/test",
				Latency:    time.Second,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// set header
			c.Writer.Header().Set("Content-Type", "application/json")

			// set status code
			c.Writer.WriteHeader(tt.status)

			event := createLogEvent(c, tt.logParams)

			if tt.expErr {
				if event == log.Info() {
					t.Errorf("createLogEvent() error = %v, expectedErr %v", event, tt.expErr)
				}
			} else {
				if event == log.Error() {
					t.Errorf("createLogEvent() = %v, want %v", event, tt.expErr)
				}
			}
		})
	}
}
