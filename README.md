# PagerDuty MCP Server (Go)

A Go implementation of the PagerDuty MCP (Model Context Protocol) server, ported from the official Python implementation.

## Features

- **50+ Tools**: Full coverage of PagerDuty API including incidents, services, teams, schedules, event orchestrations, and more
- **Read/Write Separation**: Write operations are disabled by default for safety
- **Single Binary**: Compiles to a single executable with no runtime dependencies
- **Cross-Platform**: Builds for Linux, macOS, and Windows

## Prerequisites

- Go 1.22 or later
- A PagerDuty User API Token

### Getting a PagerDuty API Token

1. Log in to PagerDuty
2. Navigate to **My Profile** > **User Settings** > **API Access**
3. Click **Create API User Token**
4. Copy and store the token securely

## Installation

### From Source

```bash
git clone https://github.com/jeremyproffitt/go-mcp-pagerduty.git
cd go-mcp-pagerduty
go build -o pagerduty-mcp ./cmd/pagerduty-mcp
```

### Using Go Install

```bash
go install github.com/jeremyproffitt/go-mcp-pagerduty/cmd/pagerduty-mcp@latest
```

## Configuration

Set the following environment variables:

```bash
export PAGERDUTY_USER_API_KEY="your-api-key-here"
# Optional: For EU accounts
export PAGERDUTY_API_HOST="https://api.eu.pagerduty.com"
```

Or create a `.env` file:

```env
PAGERDUTY_USER_API_KEY=your-api-key-here
PAGERDUTY_API_HOST=https://api.pagerduty.com
```

## Usage

### Basic (Read-Only Mode)

```bash
./pagerduty-mcp
```

### With Write Tools Enabled

```bash
./pagerduty-mcp --enable-write-tools
```

### HTTP Mode (For Containers/Lambda)

```bash
./pagerduty-mcp --http --host 0.0.0.0 --port 3000
```

### Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--http` | Run in HTTP mode | `false` |
| `--host` | HTTP server host | `127.0.0.1` |
| `--port` | HTTP server port | `3000` |
| `--enable-write-tools` | Enable write operations | `false` |

### HTTP Mode Details

When running in HTTP mode, the server exposes:
- `POST /` - MCP JSON-RPC endpoint
- `GET /health` - Health check endpoint (returns `{"status":"ok","version":"X.X.X"}`)

**Authentication**: HTTP mode requires an `Authorization` header on all requests (except `/health`). The authorization layer is pluggable; by default it accepts any token.

**Per-Request Credentials**: In HTTP mode, PagerDuty tokens can be passed via headers instead of environment variables, enabling multi-user scenarios:

| Header | Description |
|--------|-------------|
| `X-PagerDuty-Token` | PagerDuty user API key (overrides `PAGERDUTY_USER_API_KEY`) |

This header overrides the corresponding environment variable when present.

## MCP Client Configuration

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "pagerduty": {
      "command": "/path/to/pagerduty-mcp",
      "args": ["--enable-write-tools"],
      "env": {
        "PAGERDUTY_USER_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### VS Code

Add to your `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "pagerduty": {
        "type": "stdio",
        "command": "/path/to/pagerduty-mcp",
        "args": ["--enable-write-tools"],
        "env": {
          "PAGERDUTY_USER_API_KEY": "${input:pagerduty-api-key}"
        }
      }
    }
  }
}
```

### Cursor

Add to your Cursor settings:

```json
{
  "mcpServers": {
    "pagerduty": {
      "type": "stdio",
      "command": "/path/to/pagerduty-mcp",
      "args": ["--enable-write-tools"],
      "env": {
        "PAGERDUTY_USER_API_KEY": "${input:pagerduty-api-key}"
      }
    }
  }
}
```

## Docker

### Build

```bash
docker build -t pagerduty-mcp:latest .
```

### Run

```bash
docker run -i --rm \
  -e PAGERDUTY_USER_API_KEY="your-api-key-here" \
  pagerduty-mcp:latest --enable-write-tools
```

## Tool Reference

### Incidents

Tools for managing and investigating PagerDuty incidents.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_incidents` | List incidents with filtering by status, date, urgency, services, teams | `statuses`, `date_range`, `since`, `until`, `service_ids` |
| `get_incident` | Get detailed incident information by ID | `incident_id` (required) |
| `get_outlier_incident` | ML-based analysis of whether an incident is unusual | `incident_id` (required), `since` |
| `get_past_incidents` | Find similar historical incidents for troubleshooting | `incident_id` (required), `limit` |
| `get_related_incidents` | Find concurrent incidents that may be related | `incident_id` (required) |
| `list_incident_notes` | List investigation notes and comments on an incident | `incident_id` (required) |
| `create_incident` | Create a new incident manually (write) | `title`, `service_id` (required) |
| `manage_incidents` | Bulk update incidents (acknowledge, resolve, reassign) (write) | `incident_ids` (required), `status`, `urgency` |
| `add_responders` | Request additional responders for an incident (write) | `incident_id`, `responder_ids` (required) |
| `add_note_to_incident` | Add investigation note to an incident (write) | `incident_id`, `note` (required) |

### Services

Tools for managing monitored applications and components.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_services` | List services (monitored applications) | `query`, `team_ids`, `limit` |
| `get_service` | Get detailed service information | `service_id` (required) |
| `create_service` | Create a new service (write) | `name`, `escalation_policy_id` (required) |
| `update_service` | Update service configuration (write) | `service_id` (required), `name`, `description` |

### Teams

Tools for managing organizational units.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_teams` | List teams in PagerDuty | `query`, `limit` |
| `get_team` | Get team details | `team_id` (required) |
| `list_team_members` | List users in a team with their roles | `team_id` (required), `limit` |
| `create_team` | Create a new team (write) | `name` (required), `description` |
| `update_team` | Update team name or description (write) | `team_id` (required) |
| `delete_team` | DESTRUCTIVE: Delete a team permanently (write) | `team_id` (required) |
| `add_team_member` | Add user to team with role (write) | `team_id`, `user_id` (required), `role` |
| `remove_team_member` | DESTRUCTIVE: Remove user from team (write) | `team_id`, `user_id` (required) |

### Users

Tools for managing user information.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `get_user_data` | Get current authenticated user's information | None |
| `list_users` | List users in the account | `query`, `team_ids`, `limit` |

### Schedules

Tools for managing on-call rotations.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_schedules` | List on-call schedules | `query`, `limit` |
| `get_schedule` | Get schedule details with rendered on-call periods | `schedule_id` (required), `since`, `until` |
| `list_schedule_users` | List users in a schedule's rotation | `schedule_id` (required), `since`, `until` |
| `create_schedule` | Create a new schedule (write) | `name`, `time_zone` (required) |
| `create_schedule_override` | Create temporary on-call override (write) | `schedule_id`, `user_id`, `start`, `end` (required) |
| `update_schedule` | Update schedule metadata (write) | `schedule_id` (required) |

### On-Calls

Tools for finding who is currently on-call.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_oncalls` | List current and upcoming on-call entries | `earliest`, `schedule_ids`, `user_ids`, `escalation_policy_ids` |

### Escalation Policies

Tools for managing notification chains.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_escalation_policies` | List escalation policies | `query`, `user_ids`, `team_ids`, `sort_by` |
| `get_escalation_policy` | Get policy details with all levels and targets | `escalation_policy_id` (required) |

### Event Orchestrations

Tools for managing event routing and processing rules.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_event_orchestrations` | List event orchestrations (Event Rules) | `limit` |
| `get_event_orchestration` | Get orchestration details and integration URL | `orchestration_id` (required) |
| `get_event_orchestration_router` | Get router rules for service routing | `orchestration_id` (required) |
| `get_event_orchestration_global` | Get global rules (suppress, dedupe, transform) | `orchestration_id` (required) |
| `get_event_orchestration_service` | Get service-level orchestration rules | `service_id` (required) |
| `update_event_orchestration_router` | Replace entire router configuration (write) | `orchestration_id`, `config` (required) |
| `append_event_orchestration_router_rule` | Add single routing rule safely (write) | `orchestration_id`, `route_to` (required) |

### Incident Workflows

Tools for automated incident response actions.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_incident_workflows` | List available automated workflows | `query`, `limit` |
| `get_incident_workflow` | Get workflow details and configured actions | `workflow_id` (required) |
| `start_incident_workflow` | Manually trigger a workflow on an incident (write) | `workflow_id`, `incident_id` (required) |

### Change Events

Tools for correlating deployments with incidents.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_change_events` | List deployments and config changes | `since`, `until`, `team_ids`, `service_ids` |
| `get_change_event` | Get change event details | `change_event_id` (required) |
| `list_service_change_events` | List changes for a specific service | `service_id` (required), `since`, `until` |
| `list_incident_change_events` | List changes correlated with an incident | `incident_id` (required) |

### Alert Grouping

Tools for configuring how alerts combine into incidents.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_alert_grouping_settings` | List alert grouping configurations | `service_ids`, `limit` |
| `get_alert_grouping_setting` | Get grouping setting details | `setting_id` (required) |
| `create_alert_grouping_setting` | Create new grouping configuration (write) | `name`, `service_ids`, `type` (required) |
| `update_alert_grouping_setting` | Update grouping configuration (write) | `setting_id` (required) |
| `delete_alert_grouping_setting` | DESTRUCTIVE: Delete grouping setting (write) | `setting_id` (required) |

### Status Pages

Tools for public incident communication.

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `list_status_pages` | List public status pages | `limit` |
| `list_status_page_severities` | List severity options for a status page | `status_page_id` (required) |
| `list_status_page_impacts` | List impact level options | `status_page_id` (required) |
| `list_status_page_statuses` | List status lifecycle options | `status_page_id` (required) |
| `get_status_page_post` | Get incident/maintenance post details | `status_page_id`, `post_id` (required) |
| `list_status_page_post_updates` | List timeline entries for a post | `status_page_id`, `post_id` (required) |
| `create_status_page_post` | Create public incident announcement (write) | `status_page_id`, `post_type`, `title` (required) |
| `create_status_page_post_update` | Add update to existing post (write) | `status_page_id`, `post_id`, `message` (required) |

## Parameter Formats

### ID Formats

All PagerDuty IDs follow specific patterns:

| Entity | Prefix | Example |
|--------|--------|---------|
| Incident | P | `PABC123` |
| Service | P | `PDSVC123` |
| User | P | `PUSER123` |
| Team | P | `PTEAM123` |
| Schedule | P | `PSCHED123` |
| Escalation Policy | P | `PESCPOL123` |
| Workflow | P | `PWFLOW123` |
| Orchestration | E | `E1A2B3C` |

### Comma-Separated IDs

Multiple IDs can be passed as comma-separated strings:

```
service_ids: "PDSVC1,PDSVC2,PDSVC3"
team_ids: "PTEAM1,PTEAM2"
user_ids: "PUSER1,PUSER2"
incident_ids: "PABC123,PDEF456"
```

### Date/Time Formats

All date/time parameters use ISO 8601 format:

```
since: "2024-01-15T00:00:00Z"
until: "2024-01-15T23:59:59Z"
start: "2024-01-15T09:00:00-05:00"
```

### Time Zones

Use IANA time zone identifiers:

```
time_zone: "America/New_York"
time_zone: "Europe/London"
time_zone: "UTC"
time_zone: "Asia/Tokyo"
```

### Status Values

| Incident Status | Description |
|-----------------|-------------|
| `triggered` | New incident, awaiting acknowledgment |
| `acknowledged` | Responder is working on it |
| `resolved` | Incident has been fixed |

| Urgency | Description |
|---------|-------------|
| `high` | Immediate notification |
| `low` | Can wait for business hours |

### Team Roles

```
role: "manager"     # Full team management permissions
role: "responder"   # Can respond to incidents
role: "observer"    # Read-only access
```

### Alert Grouping Types

```
type: "time"           # Group alerts within time window
type: "intelligent"    # ML-based grouping
type: "content_based"  # Group by matching fields
```

## Common Workflows

### Investigating an Active Incident

1. **Find the incident**: Use `list_incidents` with `statuses: "triggered,acknowledged"`
2. **Get details**: Use `get_incident` with the incident ID
3. **Check for related changes**: Use `list_incident_change_events` to see if a recent deployment caused it
4. **Find similar past incidents**: Use `get_past_incidents` for troubleshooting guidance
5. **Check related incidents**: Use `get_related_incidents` to see if this is part of a larger issue
6. **Review notes**: Use `list_incident_notes` to see investigation progress

### Finding Who Is On-Call

1. **Current on-call**: Use `list_oncalls` with `earliest: true` to get current on-call person per schedule
2. **For specific team**: First use `list_teams` to find team ID, then `list_escalation_policies` filtered by team
3. **For specific service**: Use `get_service` to find its escalation policy, then `get_escalation_policy` for details

### Responding to an Incident

1. **Acknowledge**: Use `manage_incidents` with `status: "acknowledged"` and your incident IDs
2. **Add notes**: Use `add_note_to_incident` to document your investigation
3. **Request help**: Use `add_responders` to bring in additional team members
4. **Resolve**: Use `manage_incidents` with `status: "resolved"` when fixed

### Creating a Schedule Override (Vacation Coverage)

1. **Find the schedule**: Use `list_schedules` with a name filter
2. **Find the covering user**: Use `list_users` to get their ID
3. **Create override**: Use `create_schedule_override` with schedule ID, user ID, and time range

### Checking Service Health

1. **List all services**: Use `list_services` to see all monitored components
2. **Get service details**: Use `get_service` to see escalation policy and integrations
3. **Check recent incidents**: Use `list_incidents` with `service_ids` filter
4. **Check recent changes**: Use `list_service_change_events` to see deployments

### Setting Up Event Routing

1. **List orchestrations**: Use `list_event_orchestrations` to find existing event rules
2. **View current routing**: Use `get_event_orchestration_router` to see rules
3. **Add new rule**: Use `append_event_orchestration_router_rule` to safely add without affecting existing rules

### Communicating Incidents Publicly

1. **Find status page**: Use `list_status_pages` to get your public status page
2. **Get valid values**: Use `list_status_page_statuses` and `list_status_page_severities`
3. **Create post**: Use `create_status_page_post` with appropriate type, title, status, severity
4. **Add updates**: Use `create_status_page_post_update` as the incident progresses

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `401 Unauthorized` | Invalid or expired API token | Regenerate API token in PagerDuty |
| `403 Forbidden` | Token lacks required permissions | Check user role and API token scope |
| `404 Not Found` | Resource ID doesn't exist | Verify ID format and existence |
| `400 Bad Request` | Invalid parameter format | Check date formats, ID formats |
| `429 Too Many Requests` | Rate limit exceeded | Wait and retry with exponential backoff |

### Rate Limits

PagerDuty enforces API rate limits. If you receive a 429 error:
- Wait at least 1 second before retrying
- Use pagination (`limit` parameter) to reduce result sizes
- Cache frequently accessed data when possible

### Write Tool Errors

Write operations may fail with additional validation errors:
- **Service requires escalation policy**: Always provide `escalation_policy_id` when creating services
- **Invalid status transition**: Can only change incident status to `acknowledged` or `resolved`, not `triggered`
- **Override time conflict**: Schedule overrides cannot overlap with existing overrides

### Debugging Tips

1. **Verify IDs exist**: Use corresponding `get_*` or `list_*` tool before write operations
2. **Check permissions**: Use `get_user_data` to verify your token's user context
3. **Validate date ranges**: Ensure `since` is before `until` and dates are valid ISO 8601
4. **Check write mode**: Ensure server is started with `--enable-write-tools` for write operations

## Development

### Building

```bash
go build -o pagerduty-mcp ./cmd/pagerduty-mcp
```

### Testing

```bash
go test ./...
```

### Linting

```bash
go vet ./...
golangci-lint run
```

## License

Apache 2.0 - See LICENSE file for details.

## Acknowledgments

This is a Go port of the [official PagerDuty MCP Server](https://github.com/PagerDuty/pagerduty-mcp-server) (Python).
