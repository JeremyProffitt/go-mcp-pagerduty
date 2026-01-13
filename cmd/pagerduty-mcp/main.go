package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/auth"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/client"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/server"
	"github.com/joho/godotenv"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	enableWriteTools := flag.Bool("enable-write-tools", false, "Enable write operations (create, update, delete)")
	httpMode := flag.Bool("http", false, "Run in HTTP mode instead of stdio")
	host := flag.String("host", "127.0.0.1", "Host to listen on in HTTP mode")
	port := flag.Int("port", 3000, "Port to listen on in HTTP mode")
	flag.Parse()

	// Load .env file if it exists
	_ = godotenv.Load()

	// Print startup message
	fmt.Fprintln(os.Stderr, "Starting PagerDuty MCP Server...")
	if *enableWriteTools {
		fmt.Fprintln(os.Stderr, "Write tools ENABLED - be cautious with destructive operations")
	} else {
		fmt.Fprintln(os.Stderr, "Write tools DISABLED - use --enable-write-tools to enable")
	}

	// Create PagerDuty client
	pdClient, err := client.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create PagerDuty client: %v", err)
	}

	// Create MCP server
	mcpSrv := server.New(server.Config{
		EnableWriteTools: *enableWriteTools,
	}, pdClient)

	if *httpMode {
		// Run in HTTP mode
		fmt.Fprintf(os.Stderr, "Running in HTTP mode on %s:%d\n", *host, *port)
		httpServer := server.NewHTTPServer(mcpSrv, server.HTTPConfig{
			Host:       *host,
			Port:       *port,
			Authorizer: &auth.MockAuthorizer{},
		})
		if err := httpServer.RunHTTP(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		// Run the server on stdio
		if err := mcpserver.ServeStdio(mcpSrv); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
