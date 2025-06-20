package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var (
	port int
	host string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the fasthttp server",
	Long: `Start a high-performance HTTP server using the fasthttp library.
The server will handle HTTP requests and can be configured with different
ports and host addresses.

Examples:
  fasthttp-server server                    # Start on default port 8080
  fasthttp-server server --port 3000       # Start on port 3000
  fasthttp-server server -p 9000 -h 0.0.0.0 # Start on all interfaces port 9000`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	serverCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host address to bind the server to")
}

func startServer() {
	addr := host + ":" + strconv.Itoa(port)

	fmt.Printf("Starting fasthttp server on %s\n", addr)
	fmt.Printf("Server URL: http://%s\n", addr)

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create the server
	server := &fasthttp.Server{
		Handler: requestHandler,
		Name:    "fasthttp-server",
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(addr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	fmt.Println("\nShutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	// Log the request
	fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), method, path)

	// Set response headers
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Server", "fasthttp-server")

	// Handle different routes
	switch path {
	case "/":
		handleRoot(ctx)
	case "/health":
		handleHealth(ctx)
	case "/info":
		handleInfo(ctx)
	}
}

func handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	response := `{
  "message": "Welcome to fasthttp server!"
}`
	ctx.SetBodyString(response)
}

func handleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	response := `{
  "status": "healthy"
}`
	ctx.SetBodyString(response)
}

func handleInfo(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	response := fmt.Sprintf(`{
  "server": "fasthttp-server",
  "method": "%s",
  "path": "%s",
  "headers": %d,
  "timestamp": "%s",
  "remote_addr": "%s"
}`,
		string(ctx.Method()),
		string(ctx.Path()),
		ctx.Request.Header.Len(),
		time.Now().Format(time.RFC3339),
		ctx.RemoteAddr().String())
	ctx.SetBodyString(response)
}
