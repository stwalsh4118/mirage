# PBI-16: Clerk Authentication & User Management

[View in Backlog](../backlog.md#user-content-16)

## Overview

Implement comprehensive user authentication and management using Clerk, including frontend integration, backend JWT verification, webhook-based user synchronization, and user ownership of all Mirage resources.

## Problem Statement

Currently, Mirage has no authentication system. All resources (environments, services) are unowned, and there's no way to:
- Identify who created or owns resources
- Secure API endpoints
- Attribute actions to specific users for audit trails
- Manage per-user Railway API tokens securely

This creates security vulnerabilities, makes auditing impossible, and prevents multi-tenant usage of the platform.

## User Stories

- As a platform engineer, I want to sign up and log in using Clerk so I can securely access Mirage
- As a platform engineer, I want all my environments and services to be associated with my user account so I can manage my resources separately from others
- As a platform engineer, I want the API to automatically identify me from my session token so I don't have to manually pass credentials
- As a system admin, I want user data automatically synchronized from Clerk webhooks so our local database stays in sync
- As a developer, I want protected API routes that require authentication so unauthorized users cannot access or modify resources

## Technical Approach

### Frontend Integration
- Install and configure `@clerk/nextjs` package
- Wrap Next.js app with `<ClerkProvider>`
- Add Clerk UI components for sign-in, sign-up, and user profile
- Store and send JWT tokens with API requests
- Handle authentication state and redirects

### Backend JWT Verification
- Install Clerk Go SDK (`github.com/clerk/clerk-sdk-go/v2`)
- Create authentication middleware that:
  - Extracts JWT from Authorization header
  - Verifies JWT signature and claims using Clerk public keys
  - Extracts user ID (Clerk user ID) and adds to request context
  - Returns 401 for invalid/missing tokens
- Apply middleware to all protected routes

### User Model & Database
- Create `User` model with GORM:
  - `ID` (UUID, primary key)
  - `ClerkUserID` (string, unique, indexed) - The Clerk user ID
  - `Email` (string, unique, indexed)
  - `FirstName` (string, nullable)
  - `LastName` (string, nullable)
  - `ProfileImageURL` (string, nullable)
  - `Role` (enum: "user", "admin") - User's role
  - `CreatedAt` (timestamp)
  - `UpdatedAt` (timestamp)
  - `LastSeenAt` (timestamp, nullable)
  - `IsActive` (boolean, default true)

- Create `AuditLog` model for tracking user actions:
  - `ID` (UUID, primary key)
  - `UserID` (UUID, foreign key, indexed)
  - `Action` (string) - e.g., "environment.create", "service.delete"
  - `ResourceType` (string) - e.g., "environment", "service"
  - `ResourceID` (string, nullable) - ID of affected resource
  - `Metadata` (JSONB) - Additional context about the action
  - `IPAddress` (string)
  - `UserAgent` (string)
  - `Status` (enum: "success", "failure")
  - `ErrorMessage` (string, nullable)
  - `CreatedAt` (timestamp, indexed)

- Create `UserPreferences` model:
  - `UserID` (UUID, primary key, foreign key)
  - `Theme` (string, default "mirage") - UI theme preference
  - `DefaultEnvironmentType` (string, nullable) - Default env type for wizard
  - `NotificationsEnabled` (boolean, default true)
  - `EmailNotifications` (boolean, default true)
  - `CreatedAt` (timestamp)
  - `UpdatedAt` (timestamp)

### Resource Ownership
- Add `UserID` foreign key to existing models:
  - `Environment.UserID` (references User.ID)
  - `Service.UserID` (references User.ID)
  - `EnvironmentMetadata.UserID` (references User.ID)
- Add database indexes on UserID columns
- Update all create operations to associate resources with authenticated user
- Add user ownership checks to all read/update/delete operations

### Clerk Webhooks
- Create webhook handler endpoint `/api/webhooks/clerk`
- Verify webhook signatures using Clerk webhook secret
- Handle webhook events:
  - `user.created`: Create new User record in database
  - `user.updated`: Update User record with latest profile data
  - `user.deleted`: Soft delete or mark user as inactive
- Log all webhook events for debugging
- Return appropriate HTTP status codes (200 for success, 400/401 for errors)

### Role-Based Access Control (RBAC)
- Two roles: `user` and `admin`
- Default role: `user` (set on first webhook user.created)
- Admin promotion: Manual database update or admin API endpoint

**User Permissions:**
- View/create/update/delete their own environments and services
- Manage their own secrets in Vault
- View their own audit logs
- Update their own profile and preferences

**Admin Permissions:**
- All user permissions, plus:
- View all users and their resources (read-only)
- View system-wide audit logs
- Promote users to admin role
- Deactivate users (soft delete)
- View platform statistics and health metrics
- Access admin dashboard with system overview

### API Route Protection
- Create reusable auth middleware functions:
  - `RequireAuth()` - Requires valid JWT, extracts user
  - `RequireAdmin()` - Requires valid JWT + admin role
  - `OptionalAuth()` - Extracts user if JWT present, allows anonymous if not
  - `LoadUserContext()` - Fetches full User record from database after JWT verification
  - `AuditAction()` - Logs action to audit log
- Update all controller route registrations to use auth middleware
- Add helper functions to get current user from Gin context
- Wrap all mutating operations with audit logging

### Admin Management Endpoints
- `GET /api/v1/admin/users` - List all users with pagination and filters
- `GET /api/v1/admin/users/:id` - Get detailed user information
- `POST /api/v1/admin/users/:id/promote` - Promote user to admin
- `POST /api/v1/admin/users/:id/demote` - Demote admin to user
- `POST /api/v1/admin/users/:id/deactivate` - Deactivate user account
- `GET /api/v1/admin/environments` - List all environments across all users
- `GET /api/v1/admin/audit-logs` - View system-wide audit logs with filtering
- `GET /api/v1/admin/stats` - Get platform usage statistics

### User Management Endpoints
- `GET /api/v1/users/me` - Get current user profile
- `PATCH /api/v1/users/me` - Update current user profile (name, preferences)
- `GET /api/v1/users/me/preferences` - Get user preferences
- `PATCH /api/v1/users/me/preferences` - Update user preferences
- `GET /api/v1/users/me/audit-logs` - Get user's own audit logs
- `DELETE /api/v1/users/me` - Soft delete user account (requires confirmation)

## UX/UI Considerations

### Frontend Experience
- Clerk provides pre-built, customizable UI components
- Sign-in/sign-up flows are hosted by Clerk (can be customized)
- Seamless authentication state management
- User profile management via Clerk components
- Theme the Clerk components to match Mirage's desert/mirage aesthetic
- Settings page with tabs:
  - Profile (name, email, avatar)
  - Preferences (theme, defaults)
  - API Credentials (managed in PBI 17)
  - Audit Log (view own actions)
- Admin-only dashboard section:
  - User management table with search/filter
  - System statistics cards
  - Audit log viewer with advanced filters

### Developer Experience
- Clear error messages when authentication fails (401 Unauthorized)
- Automatic token refresh handled by Clerk SDK
- Easy access to current user in backend request handlers
- Consistent authentication across all API endpoints
- Automatic audit logging without manual instrumentation

### Error Handling
- 401 Unauthorized: Token missing, invalid, or expired
- 403 Forbidden: User authenticated but doesn't own resource or lacks admin role
- Clear error messages indicating authentication vs authorization vs permission failures
- Frontend redirects to sign-in page when authentication required
- Graceful handling of webhook delivery failures with retry logic

## Acceptance Criteria

1. **Frontend Clerk Integration**:
   - Clerk SDK installed and configured
   - Sign-in and sign-up pages work correctly
   - User can view and edit profile
   - JWT tokens are automatically sent with API requests
   - Authentication state persists across page refreshes

2. **Backend JWT Verification**:
   - Middleware successfully verifies Clerk JWTs
   - Invalid/expired tokens are rejected with 401
   - User ID is correctly extracted from verified tokens
   - Middleware is applied to all protected routes

3. **User Database Model**:
   - User table created with all required fields
   - Clerk user ID is stored and indexed
   - User data includes email, name, and profile image
   - Timestamps track creation and last seen

4. **Resource Ownership**:
   - All existing models have UserID foreign keys
   - Creating resources automatically associates with authenticated user
   - Users can only view/modify their own resources
   - Database migrations handle adding UserID to existing data

5. **Clerk Webhooks**:
   - Webhook endpoint receives and verifies Clerk events
   - User records are created when users sign up in Clerk
   - User records are updated when profile data changes in Clerk
   - User deletion/deactivation is handled appropriately
   - All webhook events are logged for debugging

6. **API Protection**:
   - All environment and service endpoints require authentication
   - Unauthenticated requests return 401
   - Requests for other users' resources return 403
   - Health check endpoint remains public
   - Webhook endpoint has its own verification mechanism

7. **Role-Based Access Control**:
   - Users assigned "user" role by default
   - Admins can promote/demote users
   - Admin-only endpoints return 403 for regular users
   - Role changes are logged in audit log
   - First user in system can be manually promoted to admin

8. **Audit Logging**:
   - All mutating operations (create, update, delete) are logged
   - Audit logs include user ID, action, resource, timestamp, IP, user agent
   - Users can view their own audit logs via UI
   - Admins can view system-wide audit logs
   - Failed operations are logged with error messages
   - Audit log table is indexed for efficient querying

9. **User Preferences**:
   - Users can set theme preference
   - Users can configure notification preferences
   - Preferences persist across sessions
   - Preferences API works correctly
   - Default values set on user creation

10. **Admin Dashboard**:
    - Admin users see "Admin" section in navigation
    - User management page lists all users with search/filter
    - Admin can view user details including owned resources
    - Admin can promote/demote users
    - System statistics show user count, environment count, etc.
    - Audit log viewer has date range, user, and action filters

## Dependencies

- Clerk account and application setup (publishable key, secret key, webhook secret)
- Frontend: `@clerk/nextjs` package
- Backend: `github.com/clerk/clerk-sdk-go/v2` package
- Database migration to add User table and UserID columns
- Environment variables for Clerk configuration:
  - `CLERK_PUBLISHABLE_KEY` (frontend)
  - `CLERK_SECRET_KEY` (backend)
  - `CLERK_WEBHOOK_SECRET` (webhook verification)

## Open Questions

1. **Migration Strategy**: ✅ ANSWERED - Delete all existing test data and start fresh

2. **User Roles**: ✅ ANSWERED - Implement basic roles (admin, user) as part of this PBI

3. **Multi-tenancy**: Should we support team/organization concepts?
   - **Decision**: Defer to future PBI, start with individual user ownership only
   - Teams would require: Organization model, team membership, resource sharing, team billing

4. **Webhook Reliability**: How to handle webhook delivery failures or out-of-order events?
   - **Proposal**: 
     - Implement idempotent webhook handlers (use Clerk webhook ID to dedupe)
     - Log all webhook events for manual replay if needed
     - Accept eventual consistency (out-of-order webhooks shouldn't break the system)
   - **Question**: Do we need a webhook queue/retry mechanism?

5. **Railway Token Migration**: ✅ ANSWERED - This PBI creates User table, PBI 17 handles Railway token migration to Vault

6. **First Admin User**: How to promote the first user to admin?
   - **Option A**: Automatically make first user in system an admin
   - **Option B**: Manual database update to set role='admin'
   - **Option C**: Environment variable `ADMIN_EMAILS` to auto-promote on webhook
   - **Recommendation**: Option C for flexibility

7. **Audit Log Retention**: How long should we keep audit logs?
   - **Proposal**: Keep indefinitely for now, add retention policy later (e.g., 90 days for regular users, 1 year for compliance)

8. **Rate Limiting**: Should we implement rate limiting per user?
   - **Proposal**: Defer to future PBI, start without rate limits

## Related Tasks

This PBI will be broken down into tasks covering:

**Infrastructure & Dependencies:**
- Research and document Clerk Go SDK (`clerk-go-guide.md`)
- Setup Clerk application (dev and prod environments)
- Configure Clerk webhooks pointing to Mirage backend

**Database Schema:**
- Create User, AuditLog, and UserPreferences models
- Create database migrations for new tables
- Add UserID foreign keys to Environment, Service, EnvironmentMetadata
- Create migration script to delete existing test data
- Add database indexes for performance (UserID, CreatedAt, Role, ClerkUserID)

**Backend Authentication:**
- Implement JWT verification middleware using Clerk SDK
- Create RequireAuth(), RequireAdmin(), OptionalAuth() middleware functions
- Create user context helpers (GetCurrentUser, GetCurrentUserID, IsAdmin)
- Implement webhook handler for Clerk events (user.created, user.updated, user.deleted)
- Add webhook signature verification
- Implement automatic admin promotion based on ADMIN_EMAILS env var

**Backend Authorization & RBAC:**
- Implement role-based access control checks
- Update Environment controller with ownership checks
- Update Service controller with ownership checks
- Add admin-only endpoints for user management
- Add admin-only endpoints for viewing all resources

**Audit Logging:**
- Create audit logging service/middleware
- Implement AuditAction() middleware wrapper
- Add audit logging to all create/update/delete operations
- Create audit log query/filter functions
- Add audit log retention policy (future-proofing)

**User Management API:**
- Implement GET /api/v1/users/me (current user profile)
- Implement PATCH /api/v1/users/me (update profile)
- Implement GET/PATCH /api/v1/users/me/preferences
- Implement GET /api/v1/users/me/audit-logs
- Implement DELETE /api/v1/users/me (soft delete)

**Admin Management API:**
- Implement GET /api/v1/admin/users (list all users)
- Implement GET /api/v1/admin/users/:id (user details)
- Implement POST /api/v1/admin/users/:id/promote
- Implement POST /api/v1/admin/users/:id/demote
- Implement POST /api/v1/admin/users/:id/deactivate
- Implement GET /api/v1/admin/environments (all environments)
- Implement GET /api/v1/admin/audit-logs (system-wide logs)
- Implement GET /api/v1/admin/stats (platform statistics)

**Frontend Clerk Integration:**
- Install and configure @clerk/nextjs
- Create Clerk provider wrapper for Next.js app
- Create sign-in page (/sign-in)
- Create sign-up page (/sign-up)
- Add user profile button/dropdown to header
- Implement authentication state management
- Add JWT token to all API requests (axios interceptor)
- Handle 401 responses with redirect to sign-in

**Frontend User Settings:**
- Create settings page layout with tabs
- Implement Profile tab (display Clerk profile)
- Implement Preferences tab (theme, defaults, notifications)
- Implement Audit Log tab (user's own actions)
- Create API Credentials tab (placeholder for PBI 17)

**Frontend Admin Dashboard:**
- Create admin route guard (check user role)
- Create admin navigation section
- Create user management page with table, search, filters
- Create user detail modal/page
- Implement promote/demote/deactivate actions
- Create system statistics dashboard
- Create admin audit log viewer with advanced filters

**Testing:**
- Unit tests for auth middleware
- Unit tests for RBAC functions
- Integration tests for webhook handlers
- Integration tests for user CRUD operations
- Integration tests for admin operations
- E2E tests for complete authentication flow
- E2E tests for admin user management

**Documentation:**
- Document authentication setup and configuration
- Document how to promote first admin user
- Document RBAC and permission model
- Document audit logging and retention
- Document admin dashboard features
- Update API documentation with auth requirements

