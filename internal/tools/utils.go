package tools

import (
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// getArgs extracts the arguments map from the request
func getArgs(request mcp.CallToolRequest) map[string]any {
	if args, ok := request.Params.Arguments.(map[string]any); ok {
		return args
	}
	return make(map[string]any)
}

// getString extracts a string argument
func getString(args map[string]any, key string) (string, bool) {
	if v, ok := args[key].(string); ok && v != "" {
		return v, true
	}
	return "", false
}

// getNumber extracts a numeric argument
func getNumber(args map[string]any, key string) (float64, bool) {
	if v, ok := args[key].(float64); ok {
		return v, true
	}
	return 0, false
}

// getBool extracts a boolean argument
func getBool(args map[string]any, key string) (bool, bool) {
	if v, ok := args[key].(bool); ok {
		return v, true
	}
	return false, false
}

// splitAndTrim splits a comma-separated string and trims whitespace
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
