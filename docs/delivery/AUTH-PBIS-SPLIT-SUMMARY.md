# Authentication PBIs Split Summary

## Overview

Based on user feedback to start with simpler, more manageable PBIs, the original comprehensive authentication design has been split into two separate PBIs:

- **PBI 16**: Core Authentication & Resource Ownership (foundational)
- **PBI 18**: Advanced Auth Features (RBAC, audit, admin dashboard)

This allows us to get basic authentication up and running quickly, then layer on advanced features once the foundation is solid.

---

## PBI 16: Core Authentication & Resource Ownership

**Status**: Proposed  
**Estimated Duration**: 2-3 weeks  
**Complexity**: Medium  

### What's Included

✅ **Frontend Clerk Integration**
- Install @clerk/nextjs
- Sign-in and sign-up pages
- User profile dropdown in header
- JWT token management
- Authentication redirects

✅ **Backend JWT Verification**
- Clerk Go SDK integration
- JWT verification middleware
- User context helpers
- Protected API routes

✅ **User Database Model**
- Basic User model (ID, ClerkUserID, Email, Name, Avatar, timestamps)
- Webhook synchronization with Clerk
- Soft delete support (IsActive flag)

✅ **Resource Ownership**
- UserID foreign keys on Environment, Service, EnvironmentMetadata
- Automatic association on resource creation
- Ownership checks on all operations (users can only access their own resources)
- 403 Forbidden for unauthorized access

✅ **User Profile API**
- GET /api/v1/users/me - View profile
- PATCH /api/v1/users/me - Update name
- GET /api/v1/users/me/environments - List own environments
- GET /api/v1/users/me/services - List own services

### What's NOT Included (Deferred to PBI 18)

❌ Role-based access control (all users equal)  
❌ Audit logging  
❌ User preferences  
❌ Admin dashboard  
❌ Admin APIs  
❌ System statistics  

### Dependencies

- Clerk account and application setup
- @clerk/nextjs (frontend)
- github.com/clerk/clerk-sdk-go/v2 (backend)
- Database migration to add User table and UserID columns

### Key Tasks (~20-25 tasks)

1. **Setup & Research** (3 tasks)
   - Research Clerk Go SDK
   - Setup Clerk application (dev + prod)
   - Configure webhooks

2. **Database** (5 tasks)
   - Create User model
   - Add UserID foreign keys to existing models
   - Create migrations
   - Delete existing test data
   - Add indexes

3. **Backend Auth** (6 tasks)
   - Install Clerk SDK
   - Implement JWT verification middleware
   - Create context helpers
   - Implement webhook handler
   - Update controllers with ownership checks
   - Implement user profile API

4. **Frontend** (6 tasks)
   - Install @clerk/nextjs
   - Setup Clerk provider
   - Create sign-in/sign-up pages
   - Add user dropdown to header
   - Add JWT to API requests
   - Handle auth redirects

5. **Testing & Docs** (5 tasks)
   - Unit tests
   - Integration tests
   - E2E tests
   - Documentation
   - Deployment

### Success Criteria

When PBI 16 is complete:
- ✅ Users can sign up and log in via Clerk
- ✅ All API endpoints require authentication
- ✅ Users can only see/manage their own resources
- ✅ User profile information is synchronized from Clerk
- ✅ Existing test data is deleted, all resources have owners

---

## PBI 18: Advanced Authentication Features

**Status**: Proposed  
**Estimated Duration**: 3-4 weeks  
**Complexity**: Medium-High  
**Depends On**: PBI 16 (must be completed first)

### What's Included

✅ **Role-Based Access Control**
- User and Admin roles
- Admin promotion via ADMIN_EMAILS env var or API
- RequireAdmin() middleware
- Role-based permission checks

✅ **Audit Logging**
- AuditLog model tracking all actions
- Automatic logging middleware
- Success/failure tracking
- Metadata storage (IP, user agent, changes)
- User and admin audit log viewers

✅ **User Preferences**
- UserPreferences model
- Theme selection
- Default environment type
- Notification settings
- Preferences API and UI

✅ **Admin Dashboard**
- User management page (list, search, filter)
- User detail view with owned resources
- Promote/demote/deactivate actions
- System statistics dashboard
- System-wide audit log viewer

✅ **Admin APIs**
- GET /api/v1/admin/users - List all users
- GET /api/v1/admin/users/:id - User details
- POST /api/v1/admin/users/:id/promote
- POST /api/v1/admin/users/:id/demote
- POST /api/v1/admin/users/:id/deactivate
- GET /api/v1/admin/environments - All environments
- GET /api/v1/admin/services - All services
- GET /api/v1/admin/stats - Platform statistics
- GET /api/v1/admin/audit-logs - System-wide logs

✅ **User Audit & Preferences APIs**
- GET /api/v1/users/me/audit-logs - User's own logs
- GET/PATCH /api/v1/users/me/preferences

### Dependencies

- PBI 16 must be completed (needs User model, auth middleware)
- No new external packages

### Key Tasks (~35-40 tasks)

1. **Database** (3 tasks)
   - Add Role field to User model
   - Create AuditLog model
   - Create UserPreferences model

2. **Backend RBAC** (4 tasks)
   - Implement RequireAdmin middleware
   - Add role to User model
   - Implement ADMIN_EMAILS auto-promotion
   - Add role checks to endpoints

3. **Backend Audit Logging** (5 tasks)
   - Create audit logging service
   - Implement AuditAction middleware
   - Add logging to all operations
   - Create query functions
   - Add audit endpoints

4. **Backend Preferences** (3 tasks)
   - Create preferences on user signup
   - Implement GET preferences endpoint
   - Implement PATCH preferences endpoint

5. **Backend Admin APIs** (9 tasks)
   - Implement all 9 admin endpoints listed above

6. **Frontend Admin Dashboard** (8 tasks)
   - Admin route guard
   - Admin navigation
   - User management page
   - User detail modal
   - Action implementations
   - Statistics dashboard
   - Audit log viewer
   - Search/filter UI

7. **Frontend User Features** (3 tasks)
   - Preferences tab in settings
   - Audit log tab in settings
   - Role badge display

8. **Testing & Docs** (5 tasks)
   - Unit tests
   - Integration tests
   - E2E tests
   - Documentation
   - Deployment

### Success Criteria

When PBI 18 is complete:
- ✅ Admin users have elevated privileges
- ✅ All user actions are logged in audit trail
- ✅ Users can customize preferences
- ✅ Admins can manage users via dashboard
- ✅ System statistics are visible to admins
- ✅ Audit logs are searchable and filterable

---

## PBI 17: HashiCorp Vault Secret Management

**Status**: Proposed (unchanged from original design)  
**Estimated Duration**: 6-8 weeks  
**Complexity**: High  
**Depends On**: PBI 16 (needs User IDs for per-user secrets)

This PBI remains as originally designed with comprehensive secret management features. See `docs/delivery/17/prd.md` for full details.

**Note**: PBI 17 can be started once PBI 16 is complete. PBI 18 can be developed in parallel or after PBI 17.

---

## Implementation Timeline

### Recommended Sequence

**Phase 1: Foundation (Weeks 1-3)**
→ **PBI 16: Core Auth**
- Get basic authentication working
- Establish resource ownership
- Users can sign in and manage their own resources

**Phase 2: Choose Your Path**

*Option A - Security First (Weeks 4-11):*
→ **PBI 17: Vault** → **PBI 18: Advanced Auth**
- Implement secret management early
- Add admin features once secrets are managed

*Option B - Admin Tools First (Weeks 4-11):*
→ **PBI 18: Advanced Auth** → **PBI 17: Vault**
- Build admin dashboard and audit logging first
- Add secret management after governance in place

*Option C - Parallel (Weeks 4-11):*
→ **PBI 17: Vault** + **PBI 18: Advanced Auth** (parallel teams)
- Fastest path if you have multiple developers
- PBI 17 and PBI 18 are independent after PBI 16

**My Recommendation**: Option A (Security First)
- Vault is critical for per-user Railway tokens
- Admin features are less urgent initially
- Can operate with one admin (manual DB promotion) short-term

---

## Comparison: Before vs After Split

### Before (Original PBI 16)
- **Duration**: 4-6 weeks
- **Tasks**: ~60 tasks
- **Complexity**: High
- **Risk**: Large scope, longer feedback cycle

### After (Split into PBI 16 + 18)
- **PBI 16 Duration**: 2-3 weeks
- **PBI 18 Duration**: 3-4 weeks
- **Total Duration**: 5-7 weeks (similar, but with milestone in middle)
- **Tasks**: ~25 + ~40 = ~65 tasks
- **Complexity**: Medium + Medium-High
- **Risk**: Lower - incremental delivery, earlier feedback

### Benefits of Split

✅ **Earlier Value Delivery**: Working auth in 2-3 weeks instead of 4-6  
✅ **Lower Risk**: Smaller chunks, easier to test and validate  
✅ **Better Feedback Loop**: Can use basic auth before building admin features  
✅ **Flexibility**: Can deprioritize PBI 18 if needed  
✅ **Clearer Dependencies**: PBI 17 only needs PBI 16, not PBI 18  
✅ **Easier Estimation**: Smaller scopes are more predictable  

---

## Migration Notes

### From Original Full Design

The full original design has been preserved at:
- `docs/delivery/16/prd-full-original.md`

If you want to reference the comprehensive design or merge features back together, that document is the authoritative source.

### Code Impact

No code has been written yet, so this split has no migration impact. We're just reorganizing the planning documents before implementation begins.

---

## Next Steps

1. ✅ **Review simplified PBI 16** (`docs/delivery/16/prd.md`)
2. ✅ **Review new PBI 18** (`docs/delivery/18/prd.md`)
3. **Approve or adjust the split**
4. **Break down PBI 16 into detailed task files**
5. **Setup Clerk application** (dev + prod)
6. **Begin implementation of PBI 16**

Once PBI 16 is complete and deployed:
7. **Choose next PBI** (17 for secrets, or 18 for admin features)
8. **Break down chosen PBI into tasks**
9. **Continue implementation**

---

## Questions?

Any concerns or adjustments needed to this split? The goal is to make PBI 16 small enough to complete in 2-3 weeks while still being valuable (working auth + resource ownership), then layer on advanced features in PBI 18.

Ready to start breaking down PBI 16 into individual task files?


