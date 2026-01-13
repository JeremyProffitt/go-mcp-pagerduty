package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/auth"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/client"
)

// createTestHandler creates an HTTP handler for testing without starting a real server
func createTestHandler(authorizer auth.Authorizer) http.Handler {
	// Create a minimal PagerDuty client (won't be used for these tests)
	pdClient := client.NewClient(client.Config{
		APIKey:  "test-api-key",
		APIHost: "https://api.pagerduty.com",
	})

	// Create MCP server
	mcpServer := New(Config{EnableWriteTools: false}, pdClient)

	// Create HTTP server instance
	httpServer := NewHTTPServer(mcpServer, HTTPConfig{
		Authorizer: authorizer,
	})

	// Build the handler (same logic as RunHTTP but without starting the server)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpServer.handleHealth)
	mux.HandleFunc("/", httpServer.handleJSONRPC)

	var handler http.Handler = mux
	if authorizer != nil {
		handler = auth.Middleware(authorizer)(mux)
	}

	return handler
}

// TestHTTPHealthEndpoint tests that GET /health returns 200 with proper JSON response
func TestHTTPHealthEndpoint(t *testing.T) {
	handler := createTestHandler(&auth.MockAuthorizer{})

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Send GET request to /health
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var healthResp healthResponse
	if err := json.Unmarshal(body, &healthResp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Validate response fields
	if healthResp.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", healthResp.Status)
	}

	if healthResp.Version == "" {
		t.Error("Expected version to be non-empty")
	}

	if healthResp.Version != ServerVersion {
		t.Errorf("Expected version '%s', got '%s'", ServerVersion, healthResp.Version)
	}
}

// TestHTTPAuthMiddleware_MissingHeader tests that POST / without Authorization header returns 401
func TestHTTPAuthMiddleware_MissingHeader(t *testing.T) {
	handler := createTestHandler(&auth.MockAuthorizer{})

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Send POST request without Authorization header
	reqBody := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`)
	resp, err := http.Post(ts.URL+"/", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code - should be 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Check response body contains error message
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !bytes.Contains(body, []byte("Authorization header required")) {
		t.Errorf("Expected error message about Authorization header, got: %s", string(body))
	}
}

// TestHTTPAuthMiddleware_WithHeader tests that POST / with Authorization header proceeds
func TestHTTPAuthMiddleware_WithHeader(t *testing.T) {
	handler := createTestHandler(&auth.MockAuthorizer{})

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Send POST request with Authorization header
	reqBody := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}`)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code - should be 200 OK (not 401)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

// TestHTTPMCPInitialize tests that POST / with valid JSON-RPC initialize request returns valid response
func TestHTTPMCPInitialize(t *testing.T) {
	handler := createTestHandler(&auth.MockAuthorizer{})

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Prepare initialize request
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	reqBody, err := json.Marshal(initRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Validate JSON-RPC response structure
	if response["jsonrpc"] != "2.0" {
		t.Errorf("Expected jsonrpc '2.0', got '%v'", response["jsonrpc"])
	}

	// Check that id matches
	if id, ok := response["id"].(float64); !ok || id != 1 {
		t.Errorf("Expected id 1, got '%v'", response["id"])
	}

	// Check that result exists and contains expected fields
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be an object, got '%v'", response["result"])
	}

	// Check for serverInfo in result
	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected serverInfo in result, got '%v'", result)
	}

	if serverInfo["name"] != ServerName {
		t.Errorf("Expected server name '%s', got '%v'", ServerName, serverInfo["name"])
	}

	if serverInfo["version"] != ServerVersion {
		t.Errorf("Expected server version '%s', got '%v'", ServerVersion, serverInfo["version"])
	}

	// Check for protocolVersion in result
	if _, ok := result["protocolVersion"]; !ok {
		t.Error("Expected protocolVersion in result")
	}

	// Check for capabilities in result
	if _, ok := result["capabilities"]; !ok {
		t.Error("Expected capabilities in result")
	}
}

// TestHTTPMCPToolsList tests that POST / with tools/list request returns list of tools
func TestHTTPMCPToolsList(t *testing.T) {
	handler := createTestHandler(&auth.MockAuthorizer{})

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// First, send initialize request (required before tools/list)
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	reqBody, _ := json.Marshal(initRequest)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send initialize request: %v", err)
	}
	resp.Body.Close()

	// Now send tools/list request
	toolsListRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	reqBody, err = json.Marshal(toolsListRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err = http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	resp, err = httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Validate JSON-RPC response structure
	if response["jsonrpc"] != "2.0" {
		t.Errorf("Expected jsonrpc '2.0', got '%v'", response["jsonrpc"])
	}

	// Check that id matches
	if id, ok := response["id"].(float64); !ok || id != 2 {
		t.Errorf("Expected id 2, got '%v'", response["id"])
	}

	// Check that result exists
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be an object, got '%v'", response["result"])
	}

	// Check that tools array exists
	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatalf("Expected tools to be an array, got '%v'", result["tools"])
	}

	// Check that there are some tools registered
	if len(tools) == 0 {
		t.Error("Expected at least one tool to be registered")
	}

	// Verify first tool has expected structure
	if len(tools) > 0 {
		firstTool, ok := tools[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected tool to be an object, got '%v'", tools[0])
		}

		// Check that tool has name field
		if _, ok := firstTool["name"]; !ok {
			t.Error("Expected tool to have 'name' field")
		}

		// Check that tool has inputSchema field
		if _, ok := firstTool["inputSchema"]; !ok {
			t.Error("Expected tool to have 'inputSchema' field")
		}
	}

	t.Logf("Successfully retrieved %d tools", len(tools))
}
