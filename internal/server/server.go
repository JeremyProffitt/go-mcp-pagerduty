package server

import (
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/client"
	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ServerName    = "PagerDuty MCP Server"
	ServerVersion = "0.1.0"
)

const MCPServerInstructions = `# PagerDuty MCP Server

## Overview
This server provides access to PagerDuty's incident management platform. When the user asks for
information about their resources, first get the user data using get_user_data and scope any
requests using the user id.

## PagerDuty Domain Concepts

### Incidents
Incidents are the core objects representing issues that need attention. They have:
- Status: triggered (new), acknowledged (being worked on), resolved (fixed)
- Urgency: high (immediate attention) or low (can wait)
- Related concepts: notes, responders, related incidents, past incidents

### Services
Services represent applications or components being monitored. Incidents are created on services.
Each service has an escalation policy that determines who gets notified.

### Teams
Teams group users together for organizational purposes. Users can belong to multiple teams.
Services and escalation policies can be associated with teams.

### Schedules
On-call schedules define who is on-call at any given time. Schedules can have overrides
for temporary coverage changes (vacations, swaps, etc.).

### Escalation Policies
Define the order in which users/schedules are notified when an incident occurs.
Multiple escalation levels ensure incidents don't go unaddressed.

### Event Orchestrations
Rules that process incoming events and route them to appropriate services.
Includes global orchestrations and service-specific orchestrations.

## Tool Categories

### Read-Only Tools (Safe)
All list_* and get_* tools are read-only and safe to use without confirmation.

### Write Tools (Use with Caution)
- create_* tools create new resources
- update_* tools modify existing resources
- manage_incidents can change incident status, urgency, and assignments
- add_* tools add relationships (responders, team members, notes)

### Destructive Tools (REQUIRES USER CONFIRMATION)
The following tools permanently delete data and should ALWAYS be confirmed with the user:
- delete_team: Permanently removes a team
- delete_alert_grouping_setting: Permanently removes an alert grouping configuration
- remove_team_member: Removes a user from a team

## Common Workflow Patterns

### Investigating an Active Incident
1. list_incidents with status=triggered,acknowledged to see active incidents
2. get_incident to see full details of a specific incident
3. list_incident_notes to see investigation notes
4. get_past_incidents to see similar historical incidents
5. get_related_incidents to see potentially related ongoing incidents
6. list_incident_change_events to see recent deployments that may have caused it

### Finding Who is On-Call
1. list_oncalls with schedule_ids or escalation_policy_ids
2. Or list_schedule_users with a date range

### Responding to an Incident
1. manage_incidents to acknowledge or resolve
2. add_note_to_incident to document findings
3. add_responders to bring in additional help

### Understanding Service Health
1. list_services to find the service
2. list_incidents filtered by service_id to see recent incidents
3. list_service_change_events to see recent deployments`

// Config holds the server configuration
type Config struct {
	EnableWriteTools bool
}

// New creates a new MCP server with the given configuration
func New(cfg Config, pdClient *client.Client) *server.MCPServer {
	s := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithInstructions(MCPServerInstructions),
	)

	// Register read-only tools (always enabled)
	registerReadTools(s, pdClient)

	// Register write tools (only if enabled)
	if cfg.EnableWriteTools {
		registerWriteTools(s, pdClient)
	}

	return s
}

// registerReadTools registers all read-only tools
func registerReadTools(s *server.MCPServer, c *client.Client) {
	// Incidents
	tools.RegisterIncidentReadTools(s, c)

	// Services
	tools.RegisterServiceReadTools(s, c)

	// Teams
	tools.RegisterTeamReadTools(s, c)

	// Users
	tools.RegisterUserReadTools(s, c)

	// Schedules
	tools.RegisterScheduleReadTools(s, c)

	// On-calls
	tools.RegisterOncallReadTools(s, c)

	// Escalation Policies
	tools.RegisterEscalationPolicyReadTools(s, c)

	// Event Orchestrations
	tools.RegisterEventOrchestrationReadTools(s, c)

	// Incident Workflows
	tools.RegisterIncidentWorkflowReadTools(s, c)

	// Change Events
	tools.RegisterChangeEventReadTools(s, c)

	// Alert Grouping Settings
	tools.RegisterAlertGroupingReadTools(s, c)

	// Status Pages
	tools.RegisterStatusPageReadTools(s, c)
}

// registerWriteTools registers all write tools
func registerWriteTools(s *server.MCPServer, c *client.Client) {
	// Incidents
	tools.RegisterIncidentWriteTools(s, c)

	// Services
	tools.RegisterServiceWriteTools(s, c)

	// Teams
	tools.RegisterTeamWriteTools(s, c)

	// Schedules
	tools.RegisterScheduleWriteTools(s, c)

	// Event Orchestrations
	tools.RegisterEventOrchestrationWriteTools(s, c)

	// Incident Workflows
	tools.RegisterIncidentWorkflowWriteTools(s, c)

	// Alert Grouping Settings
	tools.RegisterAlertGroupingWriteTools(s, c)

	// Status Pages
	tools.RegisterStatusPageWriteTools(s, c)
}
