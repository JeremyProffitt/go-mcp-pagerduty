# PagerDuty MCP Server - LLM Development Guide

This document provides guidance for LLMs working on this codebase and serves as a checklist for maintaining LLM usability standards.

## Quick Reference for Tool Selection

When a user wants to...

| User Intent | Primary Tool | Follow-up Tools |
|-------------|--------------|-----------------|
| See active incidents | `list_incidents` with `statuses: "triggered,acknowledged"` | `get_incident`, `list_incident_notes` |
| Find who is on-call | `list_oncalls` with `earliest: true` | `get_escalation_policy` |
| Investigate an incident | `get_incident`, `get_past_incidents` | `list_incident_change_events`, `get_related_incidents` |
| Respond to an incident | `manage_incidents` with `status: "acknowledged"` | `add_note_to_incident`, `add_responders` |
| See deployment impact | `list_change_events` or `list_service_change_events` | `list_incidents` for correlation |
| Find a service | `list_services` with `query` filter | `get_service` |
| Find a user | `list_users` with `query` filter | `list_oncalls` with `user_ids` |
| Create vacation coverage | `create_schedule_override` | - |
| Notify stakeholders | `create_status_page_post` | `create_status_page_post_update` |

## Tool Selection Guidance

### Similar Tools - When to Use Each

**Incident History Tools:**
- `get_past_incidents`: Use to find **similar** historical incidents (ML-based matching for troubleshooting)
- `get_related_incidents`: Use to find **concurrent** incidents (timing-based, identifies widespread issues)
- `get_outlier_incident`: Use to analyze if incident is **unusual** compared to patterns

**Change Event Tools:**
- `list_change_events`: Global view of all changes across PagerDuty
- `list_service_change_events`: Changes affecting a specific service
- `list_incident_change_events`: Changes auto-correlated with a specific incident (best for incident investigation)

**Schedule Tools:**
- `list_oncalls`: Who is on-call right now (use `earliest: true` for current only)
- `get_schedule`: Full schedule configuration and future rotation
- `list_schedule_users`: All users in a schedule's rotation

## MCP Server LLM Usability Checklist

**IMPORTANT**: This checklist must be reviewed and all items verified on every update to this repository. Any issues found must be resolved before merging changes.

### Tool Definitions

- [ ] **Clear Purpose**: Each tool has a description that clearly explains what it does and when to use it
- [ ] **No Redundant Platform Names**: Descriptions don't include unnecessary "from [Platform]" text
- [ ] **Parameter Hints**: Tool descriptions mention key parameters or capabilities
- [ ] **Use Case Guidance**: Complex tools include when-to-use hints vs similar tools
- [ ] **Consistent Naming**: All tools use snake_case naming convention
- [ ] **Action Verbs**: Tool names start with action verbs (get_, list_, create_, update_, delete_, search_)

### Parameter Documentation

- [ ] **Examples Provided**: All string parameters include format examples in descriptions
- [ ] **Format Hints**: Date/time, ID, and structured parameters document expected formats
- [ ] **Valid Values Listed**: Parameters with fixed options list valid values (e.g., "Status: 'open', 'closed', 'all'")
- [ ] **No Redundant Defaults**: Default values are in the Default field, not repeated in description text
- [ ] **Array Format Clear**: Array parameters explain expected item format
- [ ] **Object Structure Documented**: Object parameters describe expected properties

### Schema Constraints

- [ ] **Numeric Bounds**: All limit/offset/count parameters have Minimum and Maximum constraints
- [ ] **Integer Types**: Pagination and count parameters use "integer" not "number"
- [ ] **Enum Values**: Categorical parameters have Enum arrays defined in schema
- [ ] **Array Items Typed**: All array parameters have Items property with type defined
- [ ] **Object Properties**: Complex object parameters have Properties defined where structure is known
- [ ] **Pattern Validation**: ID fields have Pattern regex where format is standardized (optional)

### Tool Annotations

- [ ] **Title Set**: All tools have a human-readable Title annotation
- [ ] **ReadOnlyHint**: All get_*, list_*, search_*, describe_* tools have ReadOnlyHint: true
- [ ] **DestructiveHint**: All delete_* tools have DestructiveHint: true
- [ ] **IdempotentHint**: Safe-to-retry operations have IdempotentHint: true
- [ ] **OpenWorldHint**: Tools interacting with external systems have OpenWorldHint: true (optional)

### Token Efficiency

- [ ] **Concise Descriptions**: Tool descriptions are under 200 characters where possible
- [ ] **No Duplicate Info**: Information isn't repeated between tool and parameter descriptions
- [ ] **Abbreviated Common Terms**: Use "Max results" instead of "Maximum number of results to return"
- [ ] **Consistent Parameter Docs**: Common parameters (limit, offset, page) use identical descriptions

### Documentation

- [ ] **README Tool Reference**: README includes descriptions of what each tool does
- [ ] **Workflow Examples**: Common multi-tool workflows are documented
- [ ] **Error Handling Guide**: Common errors and recovery strategies documented
- [ ] **Parameter Patterns**: Common parameter formats (IDs, dates, queries) documented once

### Code Quality

- [ ] **Compiles Successfully**: `go build ./...` completes without errors
- [ ] **Tests Pass**: `go test ./...` completes without failures
- [ ] **No Unused Code**: No commented-out code or unused variables
- [ ] **Consistent Formatting**: Code follows Go formatting standards (`go fmt`)

## Project Structure

```
go-mcp-pagerduty/
├── cmd/pagerduty-mcp/     # Main entry point
│   └── main.go            # CLI flags, server initialization
├── internal/
│   ├── client/            # PagerDuty API client
│   │   └── client.go      # HTTP client with auth
│   ├── models/            # Data structures for API responses
│   │   ├── incidents.go
│   │   ├── services.go
│   │   └── ...
│   ├── server/            # MCP server setup
│   │   └── server.go      # Tool registration
│   └── tools/             # Tool implementations
│       ├── incidents.go   # Incident tools
│       ├── services.go    # Service tools
│       └── ...
└── README.md              # User documentation
```

## Adding New Tools

When adding a new tool:

1. **Define in appropriate tools file** (e.g., `internal/tools/incidents.go`)
2. **Follow naming convention**: `verb_noun` (e.g., `list_incidents`, `get_service`)
3. **Include comprehensive description** explaining:
   - What the tool does
   - When to use it vs similar tools
   - Key parameters to consider
4. **Add proper annotations**:
   - `mcp.WithTitleAnnotation()` for human-readable name
   - `mcp.WithReadOnlyHintAnnotation(true)` for read operations
   - `mcp.WithDestructiveHintAnnotation(true)` for delete operations
5. **Document parameters** with examples in descriptions
6. **Register in server.go** under appropriate read/write section
7. **Update README.md** Tool Reference section

## Parameter Best Practices

### ID Parameters
```go
mcp.WithString("incident_id", mcp.Required(),
    mcp.Description("The unique incident ID (e.g., 'PABC123')"))
```

### Comma-Separated IDs
```go
mcp.WithString("service_ids",
    mcp.Description("Filter by services. Comma-separated service IDs (e.g., 'PDSVC1,PDSVC2')"))
```

### Date Parameters
```go
mcp.WithString("since",
    mcp.Description("Start date in ISO 8601 format (e.g., '2024-01-15T00:00:00Z')"))
```

### Enum Parameters
```go
mcp.WithString("status",
    mcp.Description("Filter by incident status"),
    mcp.Enum("triggered", "acknowledged", "resolved"))
```

### Limit Parameters
```go
mcp.WithNumber("limit",
    mcp.Description("Maximum number of results to return"),
    mcp.Min(1), mcp.Max(100))
```

## Pre-Commit Verification

Before committing changes to this repository, run:

```bash
# Verify compilation
go build ./...

# Run all tests
go test ./...

# Check formatting
go fmt ./...

# Verify tool definitions (manual review)
# Review any new or modified tools against this checklist
```

## Issue Resolution Process

If any checklist item fails:

1. **Document the Issue**: Note which item failed and in which file
2. **Fix the Issue**: Make the necessary code changes
3. **Verify the Fix**: Re-run the relevant checks
4. **Update Tests**: Add tests for new functionality if applicable
5. **Re-verify Checklist**: Ensure fix didn't break other items

## Common Patterns

### Handler Function Structure
```go
func toolNameHandler(c *client.Client) server.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        args := getArgs(request)

        // 1. Extract and validate required parameters
        requiredParam, ok := getString(args, "required_param")
        if !ok {
            return mcp.NewToolResultError("required_param is required"), nil
        }

        // 2. Build query parameters for optional filters
        params := make(map[string]string)
        if v, ok := getString(args, "optional_param"); ok {
            params["optional_param"] = v
        }

        // 3. Make API call
        var resp models.SomeResponse
        if err := c.GetJSON("/endpoint", params, &resp); err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        // 4. Return result as JSON
        data, _ := json.Marshal(resp)
        return mcp.NewToolResultText(string(data)), nil
    }
}
```

### Error Response Format
Always return errors using `mcp.NewToolResultError()` with a clear message:
- `"parameter_name is required"` for missing required params
- `"invalid parameter_name format: expected X"` for format errors
- Pass through API errors directly for PagerDuty errors

## Testing Considerations

When testing tools:
- Use `list_*` tools before `get_*` to find valid IDs
- Test read operations before write operations
- Verify write operations with corresponding read operations
- Check error handling with invalid IDs
