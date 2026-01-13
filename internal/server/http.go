package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/auth"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// HTTPConfig holds the HTTP server configuration
type HTTPConfig struct {
	Host       string
	Port       int
	Authorizer auth.Authorizer
}

// HTTPServer wraps an MCP server with HTTP transport
type HTTPServer struct {
	mcpServer  *mcpserver.MCPServer
	config     HTTPConfig
	httpServer *http.Server
}

// NewHTTPServer creates a new HTTP server wrapping the MCP server
func NewHTTPServer(mcpServer *mcpserver.MCPServer, config HTTPConfig) *HTTPServer {
	return &HTTPServer{
		mcpServer: mcpServer,
		config:    config,
	}
}

// healthResponse represents the health check response
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// RunHTTP starts the HTTP server
func (s *HTTPServer) RunHTTP() error {
	mux := http.NewServeMux()

	// Health endpoint (no auth required)
	mux.HandleFunc("/health", s.handleHealth)

	// JSON-RPC endpoint
	mux.HandleFunc("/", s.handleJSONRPC)

	// Apply auth middleware
	var handler http.Handler = mux
	if s.config.Authorizer != nil {
		handler = auth.Middleware(s.config.Authorizer)(mux)
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Printf("Starting HTTP server on %s\n", addr)
	return s.httpServer.ListenAndServe()
}

// handleHealth handles the /health endpoint
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	resp := healthResponse{
		Status:  "ok",
		Version: ServerVersion,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleJSONRPC handles the JSON-RPC endpoint at POST /
func (s *HTTPServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Process the JSON-RPC request through the MCP server
	response := s.mcpServer.HandleMessage(r.Context(), body)

	// Marshal the response to JSON
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error":"Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
