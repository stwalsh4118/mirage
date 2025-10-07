# PBI-18: Advanced Authentication Features (RBAC, Audit, Admin Dashboard)

[View in Backlog](../backlog.md#user-content-18)

## Overview

Build advanced authentication features on top of the core Clerk authentication (PBI 16), including Role-Based Access Control, comprehensive audit logging, user preferences, and an admin dashboard for user and system management.

## Problem Statement

While PBI 16 provides core authentication and resource ownership, we need additional features for production usage:
- No way to distinguish between regular users and administrators
- No audit trail of user actions for compliance and debugging
- No way for admins to manage users or view system-wide information
- No user preferences or customization options
- Limited visibility into system usage and health

This PBI adds the governance, administration, and customization layer on top of basic authentication.

## User Stories

- As a platform admin, I want to have elevated privileges so I can manage users and view system-wide information
- As a platform admin, I want to view all users and their resources so I can provide support and oversight
- As a platform admin, I want to see an audit trail of all user actions for compliance and debugging
- As a platform admin, I want to promote/demote users to admin role as needed
- As a platform engineer, I want to view my own action history so I can track what I've done
- As a platform engineer, I want to customize my preferences (theme, defaults) so the platform works better for me
- As a system admin, I want to see platform statistics and health metrics

## Technical Approach

### Role-Based Access Control (RBAC)

**User Roles:**
- `user` (default) - Regular user with access to their own resources
- `admin` - Administrative user with elevated privileges

**Implementation:**
- Add `Role` enum field to User model
- Default all users to `user` role on creation
- Admin promotion via API endpoint or manual database update
- Environment variable `ADMIN_EMAILS` (comma-separated) for auto-promotion on signup
- Create `RequireAdmin()` middleware that checks role
- Apply admin middleware to admin-only endpoints

**Permissions:**

*User Permissions:*
- View/create/update/delete their own environments and services
- Manage their own secrets in Vault (PBI 17)
- View their own audit logs
- Update their own profile and preferences

*Admin Permissions (in addition to user permissions):*
- View all users and their profile information
- View all environments and services (read-only)
- View system-wide audit logs
- Promote users to admin role
- Demote admins to user role
- Deactivate user accounts (soft delete)
- View platform statistics and health metrics
- Access admin dashboard

**Note:** Admins cannot directly access other users' secrets (maintained in PBI 17), but can revoke/delete secrets if needed.

### Audit Logging

**AuditLog Model:**
```go
type AuditLog struct {
    ID            string    `gorm:"primaryKey;type:uuid"`
    UserID        string    `gorm:"index;not null"` // Foreign key to User
    Action        string    `gorm:"index;not null"` // e.g., "environment.create", "service.delete"
    ResourceType  string    `gorm:"index"`          // e.g., "environment", "service", "user"
    ResourceID    string    `gorm:"index"`          // ID of affected resource
    Metadata      datatypes.JSON                     // Additional context (request body, changes, etc.)
    IPAddress     string
    UserAgent     string
    Status        string    `gorm:"index"`          // "success" or "failure"
    ErrorMessage  string                             // If status=failure
    CreatedAt     time.Time `gorm:"index"`
}
```

**Implementation:**
- Create `AuditAction()` middleware wrapper
- Automatically log all mutating operations (create, update, delete)
- Log both successful and failed operations
- Store relevant metadata (old values, new values, request context)
- Query endpoints with filters:
  - By user (users see their own, admins see all)
  - By date range
  - By action type
  - By resource type
  - By status (success/failure)

**Logged Actions:**
- `environment.create`, `environment.update`, `environment.delete`
- `service.create`, `service.update`, `service.delete`
- `user.update`, `user.promote`, `user.demote`, `user.deactivate`
- `secret.*` (logged in PBI 17)
- Future: `template.create`, etc.

### User Preferences

**UserPreferences Model:**
```go
type UserPreferences struct {
    UserID                  string `gorm:"primaryKey"` // Foreign key to User
    Theme                   string `gorm:"default:'mirage'"` // UI theme
    DefaultEnvironmentType  string // Default for wizard
    NotificationsEnabled    bool   `gorm:"default:true"`
    EmailNotifications      bool   `gorm:"default:true"`
    CreatedAt               time.Time
    UpdatedAt               time.Time
}
```

**Implementation:**
- Created automatically when user signs up (webhook handler)
- API endpoints for get/update preferences
- Frontend settings page to configure
- Used in wizard to pre-select defaults

### Admin Dashboard

**User Management:**
- Table of all users with:
  - Avatar, name, email
  - Role badge
  - Created date
  - Last seen timestamp
  - Environment count
  - Service count
  - Actions dropdown
- Search by name/email
- Filter by role, active status
- Sort by various fields

**User Actions:**
- View user detail page showing:
  - Profile information
  - Owned environments (list with links)
  - Owned services (list with links)
  - Recent audit log entries
- Promote to admin
- Demote to user
- Deactivate account (soft delete, keeps resources)

**System Statistics:**
- Total users count
- Active users (last 7 days)
- Total environments
- Total services
- Environments by type (dev, staging, prod, ephemeral)
- Recent activity timeline
- Error rate (from audit logs)

**System-Wide Audit Log Viewer:**
- All user actions across the platform
- Advanced filters:
  - User selector
  - Date range picker
  - Action type multi-select
  - Resource type multi-select
  - Status (success/failure/all)
- Export to CSV
- Real-time updates (optional)

### API Endpoints

**User Management (Admin Only):**
- `GET /api/v1/admin/users` - List all users (paginated, filterable)
- `GET /api/v1/admin/users/:id` - Get user details with resources
- `POST /api/v1/admin/users/:id/promote` - Promote user to admin
- `POST /api/v1/admin/users/:id/demote` - Demote admin to user
- `POST /api/v1/admin/users/:id/deactivate` - Deactivate user account
- `GET /api/v1/admin/environments` - List all environments (all users)
- `GET /api/v1/admin/services` - List all services (all users)
- `GET /api/v1/admin/stats` - Get platform statistics
- `GET /api/v1/admin/audit-logs` - View system-wide audit logs

**User Preferences:**
- `GET /api/v1/users/me/preferences` - Get current user preferences
- `PATCH /api/v1/users/me/preferences` - Update preferences

**Audit Logs:**
- `GET /api/v1/users/me/audit-logs` - Get current user's audit logs
- Query params: `?start_date=2024-01-01&end_date=2024-12-31&action=environment.create`

## UX/UI Considerations

### Admin Navigation
- Admin users see "Admin" section in main navigation
- Admin section includes:
  - Users
  - System Stats
  - Audit Logs

### User Management Page
- Professional table with search and filters
- Inline actions (promote, demote, deactivate)
- User detail modal/page with tabs
- Clear role badges (admin vs user)
- Confirmation dialogs for destructive actions

### Settings Page
- New "Preferences" tab in user settings
- Theme selector (mirage, light, dark, etc.)
- Default environment type dropdown
- Notification toggles
- Save button with success feedback

### Audit Log Viewer
- Two views:
  - User's own audit log (in settings)
  - System-wide audit log (admin only)
- Timeline view with expandable entries
- Color coding by status (success=green, failure=red)
- Detailed metadata in expandable panels
- Export button

## Acceptance Criteria

1. **Role-Based Access Control**:
   - Users assigned "user" role by default
   - ADMIN_EMAILS env var auto-promotes specified users
   - RequireAdmin middleware blocks non-admins from admin endpoints
   - Role changes are properly validated and persisted
   - Admin role is displayed in UI

2. **Audit Logging**:
   - All mutating operations are logged automatically
   - Audit logs include user ID, action, resource, timestamp, IP, user agent
   - Both successful and failed operations are logged
   - Failed operations include error messages
   - Audit log table is indexed for efficient querying
   - Metadata is stored as JSON for flexibility

3. **User Preferences**:
   - Preferences created automatically on user signup
   - Users can view and update their preferences
   - Preferences persist across sessions
   - Theme preference is applied to UI
   - Default environment type is used in wizard

4. **Admin Dashboard**:
   - Admin users see "Admin" section in navigation
   - User management page lists all users with search/filter
   - Admin can view user details including owned resources
   - Admin can promote/demote users
   - Admin can deactivate users
   - System statistics show accurate counts
   - Audit log viewer shows all user actions with filters

5. **Admin APIs**:
   - All admin endpoints require admin role
   - Non-admin users get 403 Forbidden
   - User list endpoint supports pagination and filtering
   - User detail endpoint includes resource counts
   - Promote/demote/deactivate operations work correctly
   - Stats endpoint returns accurate system metrics

6. **User Audit Log**:
   - Users can view their own audit logs in settings
   - Audit logs are filterable by date range and action
   - Audit log entries show clear information
   - Failed actions are clearly marked
   - Audit logs are displayed in reverse chronological order

## Dependencies

- PBI 16 (Core Authentication) must be completed first
- No new external dependencies (uses existing GORM, Gin, etc.)
- Environment variable: `ADMIN_EMAILS` (comma-separated list for auto-promotion)

## Open Questions

1. **First Admin User**: How to create the first admin?
   - **Option A**: Auto-promote first user in system
   - **Option B**: Manual database update
   - **Option C**: ADMIN_EMAILS env var (recommended)
   - **Decision**: Use ADMIN_EMAILS for flexibility

2. **Audit Log Retention**: How long to keep audit logs?
   - **Proposal**: Keep indefinitely for now, add retention policy in future if needed
   - **Consideration**: Add created_at index for efficient cleanup later

3. **Deactivate vs Delete**: Should deactivated users be deleted or just marked inactive?
   - **Decision**: Soft delete (IsActive=false), keep all data for audit trail
   - **Future**: Add hard delete option for GDPR compliance if needed

4. **Audit Log Metadata**: What should we include in metadata field?
   - **Proposal**: Request body, changed fields, old values, relevant IDs
   - **Keep flexible**: JSON field allows evolution without schema changes

## Related Tasks

This PBI will be broken down into tasks covering:

**Database Schema:**
- Add Role enum field to User model
- Create AuditLog model
- Create UserPreferences model
- Create database migrations
- Add indexes for performance

**Backend RBAC:**
- Implement RequireAdmin() middleware
- Update User model with Role field
- Implement ADMIN_EMAILS auto-promotion in webhook handler
- Add role checks to admin endpoints

**Backend Audit Logging:**
- Create audit logging service
- Implement AuditAction() middleware wrapper
- Add audit logging to all mutating operations
- Create audit log query functions with filters
- Add audit log endpoints

**Backend User Preferences:**
- Implement preferences creation in webhook handler
- Implement GET /api/v1/users/me/preferences
- Implement PATCH /api/v1/users/me/preferences

**Backend Admin APIs:**
- Implement GET /api/v1/admin/users (list with pagination)
- Implement GET /api/v1/admin/users/:id (user details)
- Implement POST /api/v1/admin/users/:id/promote
- Implement POST /api/v1/admin/users/:id/demote
- Implement POST /api/v1/admin/users/:id/deactivate
- Implement GET /api/v1/admin/environments
- Implement GET /api/v1/admin/services
- Implement GET /api/v1/admin/stats
- Implement GET /api/v1/admin/audit-logs

**Frontend Admin Dashboard:**
- Create admin route guard (check role)
- Create admin navigation section
- Create user management page with table
- Implement user search and filters
- Create user detail modal/page
- Implement promote/demote/deactivate actions with confirmations
- Create system statistics dashboard
- Create admin audit log viewer

**Frontend User Preferences:**
- Add Preferences tab to settings page
- Implement theme selector
- Implement default environment type selector
- Implement notification toggles
- Implement save functionality

**Frontend User Audit Log:**
- Add Audit Log tab to settings page
- Implement audit log table/timeline
- Implement date range filter
- Implement action type filter
- Show detailed metadata on expand

**Testing:**
- Unit tests for RBAC middleware
- Unit tests for audit logging
- Integration tests for admin APIs
- Integration tests for preferences
- E2E tests for admin dashboard
- E2E tests for user management
- E2E tests for audit log viewing

**Documentation:**
- Document RBAC and permission model
- Document how to promote first admin user
- Document audit logging and retention
- Document admin dashboard features
- Update API documentation

## Notes

This PBI builds directly on PBI 16 and assumes:
- User model exists with ClerkUserID
- JWT authentication middleware is in place
- Resource ownership is implemented
- Webhook handler is receiving Clerk events

All new features are additive and don't modify PBI 16's core functionality.


