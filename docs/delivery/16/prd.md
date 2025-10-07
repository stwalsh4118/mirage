# PBI-16: Clerk Authentication & Resource Ownership (Core)

[View in Backlog](../backlog.md#user-content-16)

## Overview

Implement core user authentication using Clerk with JWT verification, user database synchronization via webhooks, and resource ownership for all Mirage entities. This is the foundational authentication layer - advanced features like RBAC, audit logging, and admin dashboards are deferred to PBI 18.

## Problem Statement

Currently, Mirage has no authentication system. All resources (environments, services) are unowned, and there's no way to:
- Identify who created or owns resources
- Secure API endpoints
- Support multiple users accessing the platform
- Manage user sessions and profiles

This PBI establishes the foundation for multi-user access with proper resource isolation.

## User Stories

- As a platform engineer, I want to sign up and log in using Clerk so I can securely access Mirage
- As a platform engineer, I want all my environments and services to be associated with my user account so I can manage my resources separately from others
- As a platform engineer, I want the API to automatically identify me from my session token so I don't have to manually pass credentials
- As a platform engineer, I want to view and update my profile information
- As a developer, I want protected API routes that require authentication so unauthorized users cannot access the platform

## Technical Approach

### Frontend Integration
- Install and configure `@clerk/nextjs` package
- Wrap Next.js app with `<ClerkProvider>`
- Add Clerk UI components for sign-in, sign-up, and user profile
- Store and send JWT tokens with API requests (axios interceptor)
- Handle authentication state and redirects

### Backend JWT Verification
- Install Clerk Go SDK (`github.com/clerk/clerk-sdk-go/v2`)
- Create authentication middleware that:
  - Extracts JWT from Authorization header
  - Verifies JWT signature and claims using Clerk public keys
  - Extracts user ID (Clerk user ID) and adds to request context
  - Returns 401 for invalid/missing tokens
- Apply middleware to all protected routes
- Create helper functions to get current user from context

### User Model & Database
- Create `User` model with GORM:
  - `ID` (UUID, primary key)
  - `ClerkUserID` (string, unique, indexed) - The Clerk user ID
  - `Email` (string, unique, indexed)
  - `FirstName` (string, nullable)
  - `LastName` (string, nullable)
  - `ProfileImageURL` (string, nullable)
  - `CreatedAt` (timestamp)
  - `UpdatedAt` (timestamp)
  - `LastSeenAt` (timestamp, nullable)
  - `IsActive` (boolean, default true)

### Resource Ownership
- Add `UserID` foreign key to existing models:
  - `Environment.UserID` (references User.ID, indexed)
  - `Service.UserID` (references User.ID, indexed)
  - `EnvironmentMetadata.UserID` (references User.ID, indexed)
- Delete existing test data (no migration needed)
- Update all create operations to associate resources with authenticated user
- Add user ownership checks to all read/update/delete operations:
  - Users can only view/modify their own resources
  - Return 403 Forbidden for accessing others' resources

### Clerk Webhooks
- Create webhook handler endpoint `/api/webhooks/clerk`
- Verify webhook signatures using Clerk webhook secret
- Handle webhook events:
  - `user.created`: Create new User record in database
  - `user.updated`: Update User record with latest profile data
  - `user.deleted`: Mark user as inactive (soft delete)
- Implement idempotent handlers (dedupe by Clerk webhook ID)
- Log all webhook events for debugging
- Return appropriate HTTP status codes (200 for success, 400/401 for errors)

### API Route Protection
- Create `RequireAuth()` middleware that:
  - Verifies JWT token
  - Loads User from database by ClerkUserID
  - Adds User to Gin context
  - Returns 401 if token invalid or 404 if user not found
- Apply to all routes except:
  - `/api/v1/healthz`
  - `/api/webhooks/clerk`
- Helper functions:
  - `GetCurrentUser(c *gin.Context) (*store.User, error)`
  - `GetCurrentUserID(c *gin.Context) (string, error)`

### User Profile Endpoints
- `GET /api/v1/users/me` - Get current user profile
  - Returns: User object with Clerk data
- `PATCH /api/v1/users/me` - Update current user profile
  - Allowed fields: FirstName, LastName (email managed by Clerk)
  - Returns: Updated user object
- `GET /api/v1/users/me/environments` - List current user's environments
- `GET /api/v1/users/me/services` - List current user's services

### Migration Strategy
- Delete all existing test data (environments, services, environment_metadata)
- Add UserID columns with NOT NULL constraint
- All new resources require authenticated user

## UX/UI Considerations

### Frontend Experience
- Clerk provides pre-built, customizable UI components
- Sign-in/sign-up flows are hosted by Clerk
- Seamless authentication state management
- User profile dropdown in header with:
  - Profile (name, email, avatar)
  - Sign out button
- Theme Clerk components to match Mirage's desert/mirage aesthetic
- Redirect to sign-in page when accessing protected routes while unauthenticated

### Error Handling
- 401 Unauthorized: Token missing, invalid, or expired → redirect to sign-in
- 403 Forbidden: User authenticated but doesn't own resource → show error message
- Clear error messages indicating authentication vs authorization failures
- Graceful handling of webhook delivery failures with retry logic

## Acceptance Criteria

1. **Frontend Clerk Integration**:
   - Clerk SDK installed and configured
   - Sign-in and sign-up pages work correctly
   - User can view profile (Clerk-provided UI)
   - JWT tokens are automatically sent with API requests
   - Authentication state persists across page refreshes
   - Unauthenticated users redirected to sign-in

2. **Backend JWT Verification**:
   - Middleware successfully verifies Clerk JWTs
   - Invalid/expired tokens are rejected with 401
   - User ID is correctly extracted from verified tokens
   - User record loaded from database and added to context
   - Middleware is applied to all protected routes
   - Helper functions work correctly

3. **User Database Model**:
   - User table created with all required fields
   - Clerk user ID is stored and indexed
   - User data includes email, name, and profile image
   - Timestamps track creation and last seen
   - IsActive flag for soft deletes

4. **Resource Ownership**:
   - All existing models have UserID foreign keys
   - Creating resources automatically associates with authenticated user
   - Users can only view/modify their own resources
   - Accessing others' resources returns 403
   - Database constraints enforce data integrity
   - Existing test data is deleted

5. **Clerk Webhooks**:
   - Webhook endpoint receives and verifies Clerk events
   - User records are created when users sign up in Clerk
   - User records are updated when profile data changes in Clerk
   - User deletion marks IsActive=false
   - Idempotent handlers prevent duplicate processing
   - All webhook events are logged for debugging

6. **API Protection**:
   - All environment and service endpoints require authentication
   - Unauthenticated requests return 401
   - Requests for other users' resources return 403
   - Health check endpoint remains public
   - Webhook endpoint has its own verification mechanism

7. **User Profile Management**:
   - GET /api/v1/users/me returns current user
   - PATCH /api/v1/users/me updates name fields
   - User can view list of their own resources
   - Profile updates reflect in UI immediately

## Dependencies

- Clerk account and application setup (publishable key, secret key, webhook secret)
- Frontend: `@clerk/nextjs` package
- Backend: `github.com/clerk/clerk-sdk-go/v2` package
- Database migration to add User table and UserID columns
- Environment variables for Clerk configuration:
  - `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY` (frontend)
  - `NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in` (frontend)
  - `NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up` (frontend)
  - `CLERK_SECRET_KEY` (backend)
  - `CLERK_WEBHOOK_SECRET` (webhook verification)

## Open Questions

1. **Webhook Reliability**: How to handle webhook delivery failures or out-of-order events?
   - **Decision**: Implement idempotent handlers (use Clerk webhook ID to dedupe), log failures for manual replay, accept eventual consistency

2. **LastSeenAt Tracking**: Should we update this on every API request?
   - **Proposal**: Update once per hour max to avoid excessive DB writes

3. **Soft Delete Behavior**: When user deleted in Clerk, should we keep their resources?
   - **Proposal**: Mark user inactive but keep resources, add "orphaned resources" handling in future PBI

## Related Tasks

This PBI will be broken down into tasks covering:

**Infrastructure & Dependencies:**
- Research and document Clerk Go SDK (`clerk-go-guide.md`)
- Setup Clerk application (dev and prod environments)
- Configure Clerk webhooks pointing to Mirage backend

**Database Schema:**
- Create User model with GORM
- Create database migration for User table
- Add UserID foreign keys to Environment, Service, EnvironmentMetadata
- Create migration script to delete existing test data
- Add database indexes (UserID, ClerkUserID, Email)

**Backend Authentication:**
- Install and configure Clerk Go SDK
- Implement JWT verification middleware (RequireAuth)
- Create user context helpers (GetCurrentUser, GetCurrentUserID)
- Implement webhook handler for Clerk events
- Add webhook signature verification
- Update LastSeenAt tracking

**Backend Authorization:**
- Update Environment controller with ownership checks
- Update Service controller with ownership checks
- Implement 403 Forbidden responses for unauthorized access
- Add UserID to all resource creation operations
- Filter list endpoints by UserID

**User Profile API:**
- Implement GET /api/v1/users/me
- Implement PATCH /api/v1/users/me
- Implement GET /api/v1/users/me/environments
- Implement GET /api/v1/users/me/services

**Frontend Clerk Integration:**
- Install and configure @clerk/nextjs
- Create Clerk provider wrapper for Next.js app
- Create sign-in page (/sign-in)
- Create sign-up page (/sign-up)
- Add user profile button/dropdown to header
- Implement authentication state management
- Add JWT token to all API requests (axios interceptor)
- Handle 401 responses with redirect to sign-in

**Frontend User Profile:**
- Add profile dropdown to header
- Display user avatar, name, email
- Link to Clerk user profile (Clerk-hosted)
- Add sign-out button

**Testing:**
- Unit tests for auth middleware
- Unit tests for user context helpers
- Integration tests for webhook handlers
- Integration tests for user CRUD operations
- Integration tests for ownership checks
- E2E tests for complete authentication flow
- E2E tests for resource isolation

**Documentation:**
- Document authentication setup and configuration
- Document Clerk application setup steps
- Document webhook configuration
- Document API authentication requirements
- Update API documentation with auth headers

## Notes for Future PBIs

The following features are intentionally deferred to **PBI 18 (Advanced Auth Features)**:
- Role-Based Access Control (admin/user roles)
- Audit logging of user actions
- User preferences and settings
- Admin dashboard and user management
- System statistics and monitoring
- Admin APIs for viewing all users/resources

This keeps PBI 16 focused on the core authentication foundation, making it easier to implement and test before adding more complex features.

