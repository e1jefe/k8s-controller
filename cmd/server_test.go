package cmd

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestServerCommand(t *testing.T) {
	// Test that server command is properly configured
	if serverCmd == nil {
		t.Fatal("serverCmd should not be nil")
	}

	if serverCmd.Use != "server" {
		t.Errorf("Expected Use to be 'server', got '%s'", serverCmd.Use)
	}

	// Test port flag exists
	portFlag := serverCmd.Flags().Lookup("port")
	if portFlag == nil {
		t.Error("Expected 'port' flag to be defined")
	}

	// Test host flag exists
	hostFlag := serverCmd.Flags().Lookup("host")
	if hostFlag == nil {
		t.Error("Expected 'host' flag to be defined")
	}
}

func TestRequestHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "health endpoint",
			path:           "/health",
			expectedStatus: 200,
		},
		{
			name:           "root endpoint",
			path:           "/",
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request context
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI(tt.path)
			ctx.Request.Header.SetMethod("GET")

			// Call the request handler
			requestHandler(ctx)

			// Check status code
			if ctx.Response.StatusCode() != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, ctx.Response.StatusCode())
			}

			// Check content type header
			contentType := string(ctx.Response.Header.Peek("Content-Type"))
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type to be 'application/json', got '%s'", contentType)
			}
		})
	}
}
