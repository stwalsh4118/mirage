# PBI 16 & 17 Expansion Summary

This document summarizes the comprehensive expansions made to both authentication and secret management PBIs based on your feedback.

## Key Decisions Incorporated

✅ **Deleting test data** - Acceptable, will start fresh  
✅ **Implementing RBAC** - Admin and user roles included in PBI 16  
✅ **Self-hosting Vault** - Detailed Railway deployment strategy in PBI 17  
✅ **Railway token was for testing** - Understood, designed per-user token model

---

## PBI 16: Clerk Authentication - Major Expansions

### New Features Added

#### 1. **Comprehensive Audit Logging**
- New `AuditLog` model tracking all user actions
- Fields: UserID, Action, ResourceType, ResourceID, Metadata (JSONB), IPAddress, UserAgent, Status, ErrorMessage
- Automatic logging middleware for all mutating operations
- User-facing audit log viewer (users see their own, admins see all)
- Searchable/filterable audit logs with date ranges

#### 2. **User Preferences System**
- New `UserPreferences` model
- Theme preferences (mirage theme system)
- Default environment types for wizard
- Notification settings (in-app and email)
- Persists across sessions

#### 3. **Role-Based Access Control (RBAC)**
- Two roles: `user` and `admin`
- Default role: `user` (auto-assigned on first signup)
- Admin promotion via API or manual database update
- Environment variable `ADMIN_EMAILS` for auto-promotion
- Clear permission model:
  - **Users**: Manage own resources, own secrets, own profile
  - **Admins**: View all users, all resources (read-only), system stats, manage roles, view all audit logs

#### 4. **Admin Dashboard**
- User management page with search/filters
- User detail view showing owned resources
- Promote/demote/deactivate user actions
- System statistics dashboard
- System-wide audit log viewer with advanced filters
- Admin navigation section

#### 5. **User Management APIs**
Complete user profile and admin management endpoints:
- `/api/v1/users/me` - Current user profile
- `/api/v1/users/me/preferences` - User preferences
- `/api/v1/users/me/audit-logs` - User's own audit trail
- `/api/v1/admin/users` - List all users (admin)
- `/api/v1/admin/users/:id/promote` - Promote to admin
- `/api/v1/admin/users/:id/demote` - Demote to user
- `/api/v1/admin/users/:id/deactivate` - Deactivate account
- `/api/v1/admin/environments` - View all environments
- `/api/v1/admin/audit-logs` - System-wide logs
- `/api/v1/admin/stats` - Platform statistics

#### 6. **Enhanced Security**
- User role checks on all admin endpoints
- Failed operation logging in audit trail
- IP address and user agent tracking
- Audit log retention planning

### Updated Acceptance Criteria
Added 5 new acceptance criteria categories:
- Role-Based Access Control (7)
- Audit Logging (8)
- User Preferences (9)
- Admin Dashboard (10)

### Comprehensive Task Breakdown
Expanded from 8 tasks to 60+ granular tasks covering:
- Database schema (5 tasks)
- Backend auth (6 tasks)
- Backend authorization/RBAC (5 tasks)
- Audit logging (5 tasks)
- User management API (5 tasks)
- Admin management API (8 tasks)
- Frontend Clerk integration (8 tasks)
- Frontend user settings (4 tasks)
- Frontend admin dashboard (7 tasks)
- Testing (7 tasks)
- Documentation (6 tasks)

---

## PBI 17: HashiCorp Vault - Major Expansions

### New Features Added

#### 1. **Comprehensive Self-Hosting Strategy**
- Detailed Railway deployment architecture
- Integrated Raft storage (no external dependencies)
- Railway volume for persistent data
- Docker Compose for local development
- Manual unseal process with migration path to cloud KMS
- Backup & disaster recovery with automated snapshots
- Monitoring and alerting strategy
- HA considerations for future scaling

#### 2. **Multiple Secret Types**
Expanded from just Railway tokens to:
- **Railway tokens** (with validation)
- **GitHub PATs** (with scope detection)
- **Docker registry credentials** (multiple registries)
- **Environment-specific secrets** (key-value pairs per environment)
- **Custom secrets** (arbitrary user secrets with tags)

#### 3. **Secret Versioning & History**
- KV v2 with version tracking (keep last 10)
- Version history viewing
- Rollback to previous versions
- Metadata per version (created_by, created_at, last_validated)

#### 4. **Secret Metadata & Organization**
- Custom metadata on all secrets
- Tags for organization
- Secret type classification
- Last validated timestamps
- Search and filter by metadata

#### 5. **Comprehensive Secret Management APIs**
60+ API endpoints across:
- **Railway token management** (5 endpoints)
- **GitHub token management** (3 endpoints)
- **Docker credentials** (4 endpoints)
- **Environment-specific secrets** (5 endpoints)
- **Generic secrets** (5 endpoints)
- **Version management** (3 endpoints)
- **Health & audit** (2 endpoints)

Each endpoint with detailed request/response specifications

#### 6. **Advanced UI Design**
- Multi-tab credentials page:
  - Railway Token tab (status, test, rotate, remove)
  - GitHub Token tab (username, scopes, validation)
  - Docker Registries tab (multiple registries, test login)
  - Custom Secrets tab (searchable, filterable, taggable)
- Version history modal with timeline view
- Environment-specific secrets tab on environment detail page
- Bulk import/export for environment secrets
- Test connection buttons for all credential types
- Live validation feedback

#### 7. **Secret Validation & Testing**
- Railway token validation via test API call
- GitHub token validation with scope detection
- Docker credential testing
- Health check endpoints
- Connection status indicators

#### 8. **Caching & Performance**
- 5-minute TTL caching layer
- Explicit cache invalidation
- Reduces Vault load without sacrificing security
- Circuit breaker for Vault failures
- Graceful degradation when Vault unavailable

#### 9. **Vault Policies & Security**
- User-secrets policy (namespaced by user ID)
- Admin policy (full access for operations)
- Backup policy (read-only for backups)
- AppRole authentication for production
- Audit logging with 30-day retention
- Sensitive data masking in logs

#### 10. **Migration & Backward Compatibility**
- Feature flag for Vault enablement
- Environment variable fallback
- Migration guide and documentation
- UI prompts for token migration
- Admin dashboard showing migration status
- Optional CLI tool for bulk migration

### Expanded Secret Storage Interface
From 8 methods to 30+ methods covering:
- Railway tokens (5 methods)
- GitHub tokens (3 methods)
- Docker credentials (4 methods)
- Environment secrets (4 methods)
- Generic secrets (5 methods)
- Version management (3 methods)
- Metadata management (2 methods)
- Health & status (2 methods)

### Updated Acceptance Criteria
Expanded from 6 to 12 comprehensive categories:
1. Vault Infrastructure
2. Vault Client Integration
3. Railway Token Management
4. GitHub Token Management
5. Docker Credentials Management
6. Environment-Specific Secrets
7. Generic Secret Management
8. Version Management
9. Secret Management UI
10. Security & Audit
11. Error Handling & Resilience
12. Migration & Backward Compatibility

### Comprehensive Task Breakdown
Expanded from 11 tasks to 100+ granular tasks covering:
- Infrastructure & dependencies (10 tasks)
- Backend Vault client (7 tasks)
- Secret storage interface (10 tasks)
- Railway client integration (7 tasks)
- API endpoints for all secret types (25 tasks)
- Database schema (4 tasks)
- Frontend credentials pages (20 tasks)
- Migration & compatibility (6 tasks)
- Testing (10 tasks)
- Security & operations (8 tasks)
- Documentation (7 tasks)

---

## Open Questions for Discussion

### PBI 16 Questions:
1. **First Admin User**: Prefer auto-promotion via `ADMIN_EMAILS` env var, manual DB update, or auto-make first user admin?
2. **Webhook Reliability**: Need webhook queue/retry mechanism or simple idempotent handlers?
3. **Audit Log Retention**: Keep indefinitely or implement retention policy?

### PBI 17 Questions:
1. **Auto-Unseal**: Start with manual unseal or invest in cloud KMS early?
2. **Cache Location**: In-memory cache OK or need Redis for multi-instance?
3. **Degraded Mode**: Should we implement database-stored encrypted secrets as last resort fallback?
4. **Backup Frequency**: Daily snapshots sufficient or also do hourly with shorter retention?
5. **Environment Sync**: Priority for auto-syncing environment secrets to Railway?

---

## Estimated Scope

### PBI 16 (Clerk Auth + RBAC + Audit):
- **Original**: ~8-10 tasks, ~2-3 weeks
- **Expanded**: ~60 tasks, ~4-6 weeks
- **Complexity**: Medium-High (Clerk integration, RBAC, comprehensive audit logging)

### PBI 17 (Vault + Multi-Secret Types):
- **Original**: ~11 tasks, ~2-3 weeks
- **Expanded**: ~100 tasks, ~6-8 weeks
- **Complexity**: High (Vault operations, multiple secret types, comprehensive UI)

### Total:
- **Combined**: ~160 tasks, ~10-14 weeks of work
- **Dependencies**: PBI 16 must complete before starting PBI 17 (need User IDs)

---

## Recommended Approach

### Phase 1: PBI 16 Foundation (Weeks 1-2)
1. Setup Clerk application
2. Implement JWT verification middleware
3. Create User, AuditLog, UserPreferences models
4. Add UserID to existing models
5. Implement webhook handler
6. Basic frontend Clerk integration

### Phase 2: PBI 16 RBAC & Admin (Weeks 3-4)
1. Implement RBAC middleware
2. Build admin APIs
3. Build user management APIs
4. Create admin dashboard frontend
5. Create user settings page
6. Implement audit logging

### Phase 3: PBI 16 Polish & Testing (Week 5)
1. Integration testing
2. E2E testing
3. Security review
4. Documentation
5. Deployment

### Phase 4: PBI 17 Infrastructure (Weeks 6-7)
1. Deploy Vault on Railway
2. Setup local Vault development
3. Implement Vault client
4. Create secret storage interface
5. Implement Railway token management

### Phase 5: PBI 17 Multi-Secret Types (Weeks 8-10)
1. Implement GitHub token management
2. Implement Docker credentials
3. Implement environment secrets
4. Implement generic secrets
5. Implement version management
6. Build all API endpoints

### Phase 6: PBI 17 Frontend & Polish (Weeks 11-13)
1. Build credentials settings page
2. Build all secret management UIs
3. Implement version history
4. Integration testing
5. E2E testing
6. Security audit

### Phase 7: Migration & Launch (Week 14)
1. Test migration path
2. Document procedures
3. Deploy to production
4. Monitor and iterate

---

## Next Steps

1. **Review expanded PRDs**: Read through both PRD documents in detail
2. **Answer open questions**: Make decisions on the questions listed above
3. **Approve scope**: Confirm the expanded scope makes sense
4. **Break down PBI 16**: Create detailed task list for PBI 16 first
5. **Setup Clerk**: Create Clerk application and get credentials
6. **Begin implementation**: Start with PBI 16 infrastructure tasks

Would you like me to proceed with breaking down PBI 16 into detailed task files?


