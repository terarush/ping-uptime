# Ping Uptime API Reference

Base URL: `/api`

All authenticated routes require `Authorization: Bearer <token>` header. Tokens obtained via `POST /api/auth/login` or `POST /api/auth/refresh`.

---

**Total endpoints: 73**

| Module | Count |
|--------|-------|
| Auth | 5 |
| Users | 6 |
| Monitors | 7 |
| Incidents | 6 |
| Notification Channels | 5 |
| Notification Logs | 1 |
| Integrations | 5 |
| Backup & Export | 5 |
| Maintenances | 5 |
| Audit Logs | 2 |
| API Tokens | 3 |
| SSL Monitors | 6 |
| Settings | 5 |
| Tags | 7 |
| Status Pages | 7 |
| Subscribers | 6 |
| Teams | 10 |
| Analytics | 3 |

---

## Table of Contents

- [Auth](#auth)
- [Users](#users)
- [Monitors](#monitors)
- [Incidents](#incidents)
- [Notification Channels](#notification-channels)
- [Notification Logs](#notification-logs)
- [Integrations](#integrations)
- [Backup & Export](#backup--export)
- [Maintenances](#maintenances)
- [Audit Logs](#audit-logs)
- [API Tokens](#api-tokens)
- [SSL Monitors](#ssl-monitors)
- [Settings](#settings)
- [Tags](#tags)
- [Status Pages](#status-pages)
- [Subscribers](#subscribers)
- [Teams](#teams)
- [Analytics](#analytics)

---

## Auth

### `GET /api/auth/setup-status`
Check if system setup is complete.

**Auth:** Public

**Response:**
```json
{
  "is_setup": false,
  "system_name": "",
  "allow_registration": true
}
```

### `POST /api/auth/setup`
Perform initial admin setup.

**Auth:** Public

**Request:**
```json
{
  "name": "Admin User",
  "email": "admin@example.com",
  "password": "password123",
  "role": "admin",
  "is_blocked": false
}
```

**Response:** Created user object (see [UserResponse](#userresponse)).

### `POST /api/auth/login`
Authenticate and receive tokens.

**Auth:** Public

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOi...",
  "refreshToken": "eyJhbGciOi...",
  "user": { ... }
}
```

### `POST /api/auth/refresh`
Refresh expired access token.

**Auth:** Public

**Request:**
```json
{
  "refreshToken": "eyJhbGciOi..."
}
```

**Response:** New `accessToken` and `refreshToken`.

### `POST /api/auth/register`
Public user registration (if enabled).

**Auth:** Public

**Request:**
```json
{
  "name": "New User",
  "email": "new@example.com",
  "password": "password123",
  "role": "member",
  "is_blocked": false
}
```

**Response:** Created user object.

---

## Users

### `GET /api/users/verify`
Verify authenticated user and return profile.

**Auth:** Auth

**Response:** [UserResponse](#userresponse).

### `GET /api/users`
List all users.

**Auth:** Auth, Admin

**Response:** `[]UserResponse`

### `GET /api/users/:id`
Get user by ID.

**Auth:** Auth, Admin

**Response:** [UserResponse](#userresponse).

### `POST /api/users`
Create new user.

**Auth:** Auth, Admin

**Request:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123",
  "role": "member",
  "is_blocked": false
}
```

**Response:** [UserResponse](#userresponse).

### `PUT /api/users/:id`
Update user.

**Auth:** Auth, Admin

**Request:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "newpassword123",
  "role": "admin",
  "is_blocked": false
}
```

`password` is optional on update (omit to keep current).

**Response:** [UserResponse](#userresponse).

### `DELETE /api/users/:id`
Delete user.

**Auth:** Auth, Admin

**Response:** Success message.

### UserResponse

```json
{
  "id": 1,
  "name": "Admin User",
  "email": "admin@example.com",
  "role": "admin",
  "is_blocked": false,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Monitors

### `GET /api/monitors`
List all monitors for authenticated user.

**Auth:** Auth

**Response:** `[]MonitorResponse`

### `GET /api/monitors/:id`
Get monitor by ID.

**Auth:** Auth

**Response:** [MonitorResponse](#monitorresponse).

### `POST /api/monitors`
Create new monitor.

**Auth:** Auth

**Request:**
```json
{
  "name": "My Website",
  "url": "https://example.com",
  "type": "http",
  "interval": 60,
  "timeout": 30,
  "check_ssl": false
}
```

`type`: `http`, `ping`, `port`, `heartbeat`. `interval` in seconds (min 1). `timeout` in seconds (min 5).

**Response:** [MonitorResponse](#monitorresponse).

### `PUT /api/monitors/:id`
Update monitor.

**Auth:** Auth

**Request:**
```json
{
  "name": "My Website",
  "url": "https://example.com",
  "type": "http",
  "interval": 60,
  "timeout": 30,
  "status": "active",
  "check_ssl": true
}
```

`status`: `active`, `paused`.

**Response:** [MonitorResponse](#monitorresponse).

### `DELETE /api/monitors/:id`
Delete monitor.

**Auth:** Auth

**Response:** Success message.

### `GET /api/monitors/events`
SSE endpoint for real-time monitor events.

**Auth:** Public

Returns Server-Sent Events stream with monitor status changes.

### `POST /api/heartbeat/:token`
Register heartbeat ping from external service.

**Auth:** Public (token-based)

Path param `:token` is the monitor's heartbeat token.

**Response:** Success message.

### MonitorResponse

```json
{
  "id": 1,
  "name": "My Website",
  "url": "https://example.com",
  "type": "http",
  "interval": 60,
  "timeout": 30,
  "status": "active",
  "uptime_status": "up",
  "last_checked_at": "2025-01-01T00:01:00Z",
  "last_latency": 120,
  "check_ssl": true,
  "ssl_expires_at": "2025-06-01T00:00:00Z",
  "heartbeat_token": "tok_abc123",
  "user_id": 1,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Incidents

### `GET /api/incidents`
List all incidents.

**Auth:** Auth

**Response:** `[]IncidentResponse`

### `GET /api/incidents/:id`
Get incident by ID.

**Auth:** Auth

**Response:** [IncidentResponse](#incidentresponse).

### `GET /api/incidents/monitor/:monitorId`
Get incidents for a specific monitor.

**Auth:** Auth

**Response:** `[]IncidentResponse`

### `POST /api/incidents`
Create incident.

**Auth:** Auth

**Request:**
```json
{
  "monitor_id": 1,
  "status": "down",
  "error_message": "Connection timeout",
  "latency": 5000
}
```

`status`: `down`, `degraded`.

**Response:** [IncidentResponse](#incidentresponse).

### `PUT /api/incidents/:id`
Update incident.

**Auth:** Auth

**Request:**
```json
{
  "status": "resolved",
  "error_message": "Connection restored",
  "latency": 200
}
```

**Response:** [IncidentResponse](#incidentresponse).

### `DELETE /api/incidents/:id`
Delete incident.

**Auth:** Auth

**Response:** Success message.

### IncidentResponse

```json
{
  "id": 1,
  "monitor_id": 1,
  "user_id": 1,
  "status": "down",
  "error_message": "Connection timeout",
  "latency": 5000,
  "created_at": "2025-01-01T00:00:00Z",
  "resolved_at": null
}
```

---

## Notification Channels

### `GET /api/notification-channels`
List all notification channels.

**Auth:** Auth

**Response:** `[]ChannelResponse`

### `GET /api/notification-channels/:id`
Get channel by ID.

**Auth:** Auth

**Response:** [ChannelResponse](#channelresponse).

### `POST /api/notification-channels`
Create notification channel.

**Auth:** Auth

**Request:**
```json
{
  "name": "Slack Alerts",
  "type": "slack",
  "config": "{\"webhook_url\": \"https://hooks.slack.com/...\"}",
  "enabled": true
}
```

`type`: `email`, `slack`, `discord`, `webhook`, `telegram`, `sms`. `config` contains JSON-stringified provider-specific settings.

**Response:** [ChannelResponse](#channelresponse).

### `PUT /api/notification-channels/:id`
Update notification channel.

**Auth:** Auth

**Request:**
```json
{
  "name": "Slack Alerts",
  "type": "slack",
  "config": "{\"webhook_url\": \"https://hooks.slack.com/...\"}",
  "enabled": false
}
```

**Response:** [ChannelResponse](#channelresponse).

### `DELETE /api/notification-channels/:id`
Delete notification channel.

**Auth:** Auth

**Response:** Success message.

### ChannelResponse

```json
{
  "id": 1,
  "name": "Slack Alerts",
  "type": "slack",
  "config": "{\"webhook_url\": \"https://hooks.slack.com/...\"}",
  "enabled": true,
  "user_id": 1,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Notification Logs

### `GET /api/notification-logs`
List notification delivery logs.

**Auth:** Auth

**Query params:** `channel_type`, `status` (`sent`, `failed`), `event_type`, `date_from`, `date_to`, `limit`, `offset`.

**Response:** Paginated array of notification log entries.

---

## Integrations

### `GET /api/integrations`
List all integrations.

**Auth:** Auth

**Response:** `[]IntegrationResponse`

### `POST /api/integrations`
Create integration.

**Auth:** Auth

**Request:**
```json
{
  "name": "PagerDuty",
  "type": "pagerduty",
  "config": "{\"routing_key\": \"...\"}"
}
```

`type`: `slack`, `discord`, `webhook`, `github`, `pagerduty`. `config` is provider-specific JSON.

**Response:** [IntegrationResponse](#integrationresponse).

### `PUT /api/integrations/:id`
Update integration.

**Auth:** Auth

**Request:**
```json
{
  "name": "PagerDuty",
  "type": "pagerduty",
  "config": "{\"routing_key\": \"...\"}",
  "enabled": true
}
```

**Response:** [IntegrationResponse](#integrationresponse).

### `DELETE /api/integrations/:id`
Delete integration.

**Auth:** Auth

**Response:** Success message.

### `POST /api/integrations/:id/test`
Test integration connection.

**Auth:** Auth

**Response:** Test result with success/failure.

### IntegrationResponse

```json
{
  "id": 1,
  "name": "PagerDuty",
  "type": "pagerduty",
  "config": "{\"routing_key\": \"...\"}",
  "enabled": true,
  "user_id": 1,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Backup & Export

All backup endpoints require **Auth + Admin**.

### `GET /api/backup/export`
Export system configuration as JSON.

**Response:** Full system configuration JSON (monitors, channels, settings, etc.).

### `GET /api/backup/download`
Download backup as JSON file (attachment).

### `POST /api/backup/import`
Import configuration from backup file.

**Request:** Multipart form with `file` field containing JSON backup.

### `GET /api/backup/history`
List backup history.

**Response:** Array of backup record metadata.

### `DELETE /api/backup/history/:id`
Delete a backup record.

---

## Maintenances

### `GET /api/maintenances`
List all maintenances.

**Auth:** Auth

**Response:** `[]MaintenanceResponse`

### `GET /api/maintenances/:id`
Get maintenance by ID.

**Auth:** Auth

**Response:** [MaintenanceResponse](#maintenanceresponse).

### `POST /api/maintenances`
Create maintenance window.

**Auth:** Auth

**Request:**
```json
{
  "name": "Database Migration",
  "description": "Scheduled maintenance for DB upgrade",
  "start_at": "2025-01-15T02:00:00Z",
  "end_at": "2025-01-15T04:00:00Z",
  "monitor_ids": [1, 2, 3]
}
```

Times in RFC3339 format.

**Response:** [MaintenanceResponse](#maintenanceresponse).

### `PUT /api/maintenances/:id`
Update maintenance.

**Auth:** Auth

**Request:** Same structure as create.

**Response:** [MaintenanceResponse](#maintenanceresponse).

### `DELETE /api/maintenances/:id`
Delete maintenance.

**Auth:** Auth

**Response:** Success message.

### MaintenanceResponse

```json
{
  "id": 1,
  "name": "Database Migration",
  "description": "Scheduled maintenance for DB upgrade",
  "start_at": "2025-01-15T02:00:00Z",
  "end_at": "2025-01-15T04:00:00Z",
  "status": "scheduled",
  "user_id": 1,
  "monitor_ids": [1, 2, 3],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

`status`: `scheduled`, `ongoing`, `completed`, `cancelled`.

---

## Audit Logs

All audit log endpoints require **Auth + Admin**.

### `GET /api/audit-logs`
List all audit log entries.

**Response:** Array of audit log records.

### `GET /api/audit-logs/search`
Search audit logs by entity.

**Query params:** `entity_type`, `entity_id`.

**Response:** Filtered array of audit log records.

---

## API Tokens

### `GET /api/api-tokens`
List all API tokens for authenticated user.

**Auth:** Auth

**Response:** `[]TokenResponse`

### `POST /api/api-tokens`
Create new API token.

**Auth:** Auth

**Request:**
```json
{
  "name": "CI Pipeline Token",
  "expires_at": "2026-01-01T00:00:00Z"
}
```

`expires_at` is optional (RFC3339). Omit for no expiry.

**Response:**
```json
{
  "id": 1,
  "name": "CI Pipeline Token",
  "token_prefix": "pt_abc",
  "raw_token": "pt_abc123...",
  "last_used_at": null,
  "expires_at": "2026-01-01T00:00:00Z",
  "is_revoked": false,
  "created_at": "2025-01-01T00:00:00Z"
}
```

`raw_token` is shown **only on creation** — save it immediately.

### `DELETE /api/api-tokens/:id`
Revoke API token.

**Auth:** Auth

**Response:** Success message.

---

## SSL Monitors

### `GET /api/ssl-monitors`
List all SSL certificate checks.

**Auth:** Auth

**Response:** `[]SSLCertificate`

### `GET /api/ssl-monitors/:id`
Get SSL certificate by ID.

**Auth:** Auth

### `POST /api/ssl-monitors/check/:monitorId`
Check SSL for specific monitor.

**Auth:** Auth

**Response:**
```json
{
  "monitor_id": 1,
  "domain": "example.com",
  "status": "valid",
  "error": ""
}
```

### `POST /api/ssl-monitors/check-all`
Check SSL for all monitors.

**Auth:** Auth

**Response:** `[]checkResult`

### `GET /api/ssl-monitors/expiring`
Get certificates expiring soon.

**Auth:** Auth

**Query param:** `days` (default 30).

**Response:** `[]SSLCertificate` filtered by expiry.

### `DELETE /api/ssl-monitors/:id`
Delete SSL certificate record.

**Auth:** Auth

---

## Settings

### `GET /api/settings/public/system-name`
Get public system name.

**Auth:** Public

**Response:**
```json
{
  "value": "Ping Uptime"
}
```

### `GET /api/settings`
List all settings.

**Auth:** Auth

**Response:** `[]SettingResponse`

### `GET /api/settings/:key`
Get setting by key.

**Auth:** Auth

**Response:** `SettingResponse`

### `POST /api/settings`
Save or update setting.

**Auth:** Auth (checks admin role in handler)

**Request:**
```json
{
  "key": "system_name",
  "value": "My Uptime Monitor",
  "description": "Display name for the system"
}
```

**Response:** [SettingResponse](#settingresponse).

### `DELETE /api/settings/:key`
Delete setting.

**Auth:** Auth (checks admin role in handler)

### SettingResponse

```json
{
  "key": "system_name",
  "value": "My Uptime Monitor",
  "description": "Display name for the system",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Tags

### `GET /api/tags`
List all tags.

**Auth:** Auth

**Response:** Array of tag objects.

### `POST /api/tags`
Create tag.

**Auth:** Auth

**Request:**
```json
{
  "name": "production",
  "color": "#ff0000"
}
```

### `PUT /api/tags/:id`
Update tag.

**Auth:** Auth

**Request:**
```json
{
  "name": "staging",
  "color": "#00ff00"
}
```

### `DELETE /api/tags/:id`
Delete tag.

**Auth:** Auth

### `POST /api/monitors/:id/tags`
Attach tags to monitor.

**Auth:** Auth

**Request:**
```json
{
  "tag_ids": [1, 2, 3]
}
```

### `GET /api/monitors/:id/tags`
Get tags for monitor.

**Auth:** Auth

### `DELETE /api/monitors/:id/tags/:tagID`
Detach tag from monitor.

**Auth:** Auth

---

## Status Pages

### `GET /api/status-pages/slug/:slug`
Get public status page by slug.

**Auth:** Public

**Response:** Status page with monitors and current status.

### `GET /api/status-pages/:slug/badge.svg`
Get SVG status badge for status page.

**Auth:** Public

**Response:** SVG image.

### `GET /api/status-pages`
List all status pages.

**Auth:** Auth

**Response:** `[]StatusPageResponse`

### `GET /api/status-pages/:id`
Get status page by ID.

**Auth:** Auth

**Response:** [StatusPageResponse](#statuspageresponse).

### `POST /api/status-pages`
Create status page.

**Auth:** Auth

**Request:**
```json
{
  "name": "System Status",
  "slug": "system-status",
  "description": "Public status page for our services",
  "monitor_ids": [1, 2, 3]
}
```

**Response:** [StatusPageResponse](#statuspageresponse).

### `PUT /api/status-pages/:id`
Update status page.

**Auth:** Auth

**Request:** Same structure as create.

**Response:** [StatusPageResponse](#statuspageresponse).

### `DELETE /api/status-pages/:id`
Delete status page.

**Auth:** Auth

### StatusPageResponse

```json
{
  "id": 1,
  "name": "System Status",
  "slug": "system-status",
  "description": "Public status page for our services",
  "user_id": 1,
  "monitors": [
    {
      "id": 1,
      "name": "Main Website",
      "uptime_status": "up",
      ...
    }
  ],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## Subscribers

### `POST /api/status-pages/subscribe`
Subscribe to status page notifications.

**Auth:** Public

**Request:**
```json
{
  "email": "user@example.com",
  "status_page_id": 1
}
```

### `GET /api/status-pages/subscribe/verify`
Verify subscription.

**Auth:** Public

**Query param:** `token`.

### `GET /api/status-pages/unsubscribe`
Unsubscribe by token.

**Auth:** Public

**Query param:** `token`.

### `GET /api/status-pages/:pageID/subscribers/count`
Count subscribers for a status page.

**Auth:** Auth

**Response:**
```json
{
  "count": 42
}
```

### `GET /api/status-pages/:pageID/subscribers`
List subscribers for a status page.

**Auth:** Auth, Admin

### `POST /api/status-pages/:pageID/subscribers/unsubscribe`
Admin-initiated unsubscribe.

**Auth:** Auth, Admin

**Request:**
```json
{
  "email": "user@example.com",
  "status_page_id": 1
}
```

---

## Teams

### `GET /api/teams`
List all teams.

**Auth:** Auth

**Response:** `[]TeamResponse`

### `POST /api/teams`
Create team.

**Auth:** Auth

**Request:**
```json
{
  "name": "DevOps"
}
```

**Response:** [TeamResponse](#teamresponse).

### `PUT /api/teams/:id`
Update team.

**Auth:** Auth

**Request:**
```json
{
  "name": "Platform Engineering"
}
```

**Response:** [TeamResponse](#teamresponse).

### `DELETE /api/teams/:id`
Delete team.

**Auth:** Auth (admin-gated in handler)

### `GET /api/teams/:id/members`
List team members.

**Auth:** Auth

**Response:** `[]TeamMemberResponse`

### `POST /api/teams/:id/members`
Invite member to team.

**Auth:** Auth

**Request:**
```json
{
  "user_id": 2,
  "role": "member"
}
```

`role`: `admin`, `member`.

**Response:** [TeamMemberResponse](#teammemberresponse).

### `PUT /api/teams/:id/members/:memberID`
Update team member role.

**Auth:** Auth

**Request:**
```json
{
  "role": "admin"
}
```

### `DELETE /api/teams/:id/members/:memberID`
Remove team member.

**Auth:** Auth

### `POST /api/teams/:id/members/accept`
Accept team invitation.

**Auth:** Auth

### `POST /api/teams/:id/members/reject`
Reject team invitation.

**Auth:** Auth

### TeamResponse

```json
{
  "id": 1,
  "name": "DevOps",
  "member_count": 5,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### TeamMemberResponse

```json
{
  "id": 1,
  "team_id": 1,
  "user_id": 2,
  "role": "member",
  "status": "pending",
  "invited_by": 1,
  "user": {
    "id": 2,
    "email": "jane@example.com",
    "name": "Jane Doe"
  },
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

`status`: `pending`, `accepted`, `rejected`.

---

## Analytics

### `GET /api/analytics/monitors/:id/chart`
Get chart data for a monitor.

**Auth:** Auth

**Query params:** `window` (default `"1m"`). Values: `"1h"`, `"24h"`, `"7d"`, `"30d"`, `"1m"`, `"3m"`, `"1y"`.

**Response:**
```json
{
  "monitor_id": 1,
  "window": "1m",
  "data": [...]
}
```

### `GET /api/analytics/dashboard`
Get dashboard statistics.

**Auth:** Auth

**Query params:** `window` (default `"1m"`).

**Response:**
```json
{
  "window": "1m",
  "data": { ... }
}
```

### `GET /api/analytics/report`
Get analytics report.

**Auth:** Auth

**Query params:** `from` (YYYY-MM-DD, required), `to` (YYYY-MM-DD, required), `monitor_id` (optional, admin-only filter).

**Response:**
```json
{
  "from": "2025-01-01",
  "to": "2025-01-31",
  "stats": { ... }
}
```

---

## Authentication Summary

| Level | Description | Middleware |
|-------|-------------|-----------|
| **Public** | No auth required | None |
| **Auth** | Valid Bearer token required | `middleware.Auth` |
| **Admin** | Auth + admin role required | `middleware.Admin` (or inline role check in settings) |

## Common Response Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient role) |
| 404 | Resource not found |
| 500 | Internal server error |
