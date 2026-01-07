package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/client"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/server"
	"github.com/joho/godotenv"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse command line flags
	enableWriteTools := flag.Bool("enable-write-tools", false, "Enable write operations (create, update, delete)")
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
	mcpServer := server.New(server.Config{
		EnableWriteTools: *enableWriteTools,
	}, pdClient)

	// Run the server on stdio
	if err := mcpserver.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
