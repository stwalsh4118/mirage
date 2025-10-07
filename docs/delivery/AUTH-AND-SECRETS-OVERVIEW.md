# Authentication & Secret Management Design Overview

This document provides a high-level overview of PBI 16 (Clerk Authentication) and PBI 17 (HashiCorp Vault), explaining how they work together to provide secure, multi-tenant access to Mirage.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (Next.js)                       │
│  ┌────────────────┐         ┌──────────────────────────────┐   │
│  │ Clerk Provider │────────▶│ Sign In / Sign Up / Profile  │   │
│  └────────────────┘         └──────────────────────────────┘   │
│         │                                                        │
│         │ JWT Token                                              │
│         ▼                                                        │
└─────────────────────────────────────────────────────────────────┘
          │
          │ Authorization: Bearer <JWT>
          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go + Gin)                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Auth Middleware                                          │   │
│  │ • Verify JWT with Clerk                                  │   │
│  │ • Extract User ID                                        │   │
│  │ • Load User from Database                                │   │
│  │ • Add User to Request Context                            │   │
│  └─────────────────────────────────────────────────────────┘   │
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Controllers (Environment, Service, Secrets)              │   │
│  │ • Get User from Context                                  │   │
│  │ • Fetch Railway Token from Vault                         │   │
│  │ • Create User-Specific Railway Client                    │   │
│  │ • Execute Railway Operations                             │   │
│  └─────────────────────────────────────────────────────────┘   │
│         │                              │                         │
└─────────┼──────────────────────────────┼─────────────────────────┘
          │                              │
          │                              └──────────────┐
          ▼                                             ▼
┌──────────────────────┐                    ┌───────────────────┐
│   PostgreSQL DB      │                    │  HashiCorp Vault  │
│  ┌────────────────┐  │                    │ ┌───────────────┐ │
│  │ Users Table    │  │                    │ │ User Secrets  │ │
│  │ • ClerkUserID  │  │                    │ │ Railway Token │ │
│  │ • Email        │  │                    │ │ GitHub Token  │ │
│  │ • Name         │  │                    │ │ Docker Creds  │ │
│  └────────────────┘  │                    │ └───────────────┘ │
│  ┌────────────────┐  │                    └───────────────────┘
│  │ Environments   │  │                              
│  │ • UserID (FK)  │  │                    
│  └────────────────┘  │                    
│  ┌────────────────┐  │                    
│  │ Services       │  │                    
│  │ • UserID (FK)  │  │                    
│  └────────────────┘  │                    
└──────────────────────┘                    

         ▲
         │ Webhook Events
         │ (user.created, user.updated, user.deleted)
         │
┌──────────────────────┐
│   Clerk Service      │
│  • User Management   │
│  • JWT Issuance      │
│  • Profile Data      │
└──────────────────────┘
```

## Key Design Decisions

### 1. Clerk for Authentication (PBI 16)

**Why Clerk?**
- ✅ Simple frontend integration with pre-built React components
- ✅ Robust Go SDK for JWT verification
- ✅ Webhooks for automatic user synchronization
- ✅ Handles password management, MFA, OAuth providers
- ✅ Professional UI that can be themed to match Mirage
- ✅ Free tier sufficient for early stage

**Alternative Considered**: Auth0, Firebase Auth, Custom JWT solution
**Decision**: Clerk provides the best balance of ease-of-use and features

### 2. HashiCorp Vault for Secrets (PBI 17)

**Why Vault?**
- ✅ Industry-standard secret management solution
- ✅ Encryption at rest and in transit
- ✅ Secret versioning for rotation support
- ✅ Audit logging of all secret access
- ✅ Flexible auth methods (token, AppRole, cloud IAM)
- ✅ Can use HCP Vault (managed) or self-host

**Alternative Considered**: Cloud provider secret managers (AWS Secrets Manager, GCP Secret Manager), Database encryption
**Decision**: Vault provides best control and flexibility across cloud providers

### 3. Per-User Railway Tokens

**Current State**: Single `RAILWAY_API_TOKEN` environment variable shared by all users

**New State**: Each user stores their own Railway token in Vault
- ✅ Actions attributed to correct Railway account
- ✅ Users manage their own Railway access
- ✅ Individual token revocation without affecting others
- ✅ Supports different Railway projects per user

**Migration Path**: 
1. Deploy with Vault optional (env var fallback)
2. Prompt users to configure tokens
3. Eventually require Vault tokens

## Implementation Sequence

### Phase 1: PBI 16 - Authentication Foundation
1. Frontend Clerk integration
2. Backend JWT verification middleware
3. User database model and migrations
4. Webhook handler for user sync
5. Resource ownership (UserID foreign keys)
6. API route protection

**Outcome**: Users can sign in, API is secured, resources are owned

### Phase 2: PBI 17 - Secret Management
1. Vault server setup (HCP Vault for prod, Docker for dev)
2. Vault Go client integration
3. Secret storage interface and implementation
4. Railway token management API
5. Update Railway client to use per-user tokens
6. Frontend UI for token configuration
7. Migration from env var tokens

**Outcome**: Users manage their own Railway tokens securely

## Data Flow Examples

### Example 1: User Signs Up
```
1. User clicks "Sign Up" in Mirage frontend
2. Clerk handles sign-up flow (email/password, OAuth, etc.)
3. Clerk sends webhook to Mirage: POST /api/webhooks/clerk
4. Webhook handler creates User record in PostgreSQL
5. User is redirected to Mirage dashboard
6. Dashboard shows "Configure Railway Token" prompt
```

### Example 2: Creating an Environment
```
1. User fills out environment creation form
2. Frontend sends POST /api/v1/environments with JWT token
3. Auth middleware verifies JWT with Clerk, extracts user ID
4. Controller gets user from context
5. Controller fetches user's Railway token from Vault
6. Controller creates Railway client with user's token
7. Controller calls Railway API to create project/environment
8. Environment record created in DB with UserID
9. Response returned to frontend
```

### Example 3: Storing Railway Token
```
1. User enters Railway token in settings page
2. Frontend sends POST /api/v1/secrets/railway
3. Auth middleware verifies user
4. Controller validates token by testing Railway API
5. Controller stores token in Vault at /mirage/users/{userID}/railway
6. Success response returned
7. User can now create environments
```

## Security Considerations

### Authentication Security
- JWTs verified using Clerk's public keys (JWKS)
- Tokens expire and must be refreshed
- Webhook signatures verified to prevent spoofing
- User ID extracted from verified JWT, not from request body

### Secret Security
- Secrets never logged in plain text
- Secrets encrypted at rest in Vault
- TLS for all Vault communication
- Secrets cached briefly (5 min) to reduce Vault load
- Audit log tracks all secret access

### Authorization Security
- Users can only access their own resources
- Foreign key constraints enforce data integrity
- Database queries filtered by UserID
- 403 Forbidden for accessing others' resources

## Configuration

### Environment Variables - PBI 16 (Clerk)
```bash
# Frontend (.env.local)
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in
NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up

# Backend
CLERK_SECRET_KEY=sk_test_...
CLERK_WEBHOOK_SECRET=whsec_...
```

### Environment Variables - PBI 17 (Vault)
```bash
# Backend
VAULT_ADDR=https://vault.example.com:8200
VAULT_TOKEN=hvs.CAESI...  # or AppRole credentials
VAULT_NAMESPACE=mirage     # Enterprise only
VAULT_SKIP_VERIFY=false    # true for local dev only
```

## Open Questions for Discussion

### PBI 16 Questions:
1. **Existing Data Migration**: What should we do with existing environments/services that have no owner?
   - My recommendation: Since this is early stage, we can delete existing test data and start fresh

2. **User Roles**: Should we implement admin/user roles now?
   - My recommendation: Defer to later PBI, start with equal permissions for all authenticated users

3. **Webhook Reliability**: How to handle webhook failures?
   - My recommendation: Basic retry with idempotent handlers, log failures for manual review

### PBI 17 Questions:
1. **Vault Hosting**: Should we use HCP Vault (managed) or self-host?
   - My recommendation: Start with HCP Vault for simplicity, offers free tier

2. **Secret Caching**: How long should we cache Railway tokens?
   - My recommendation: 5 minutes, invalidate on update/delete

3. **Migration Timeline**: When should we stop supporting env var tokens?
   - My recommendation: Support for 1-2 release cycles, then deprecate

4. **Admin Secret Access**: Should admins be able to view users' secrets?
   - My recommendation: No direct access, but ability to revoke if needed

## Next Steps

1. **Review PRDs**: Read through PBI 16 and PBI 17 PRD documents
2. **Answer Open Questions**: Decide on open questions above
3. **Break Down into Tasks**: Create detailed task lists for each PBI
4. **Research External Packages**: Create package guides for Clerk SDK and Vault SDK
5. **Set Up Accounts**: Create Clerk application and HCP Vault cluster
6. **Begin Implementation**: Start with PBI 16 (auth foundation required for PBI 17)

## References

- [PBI 16 PRD - Clerk Authentication](./16/prd.md)
- [PBI 17 PRD - HashiCorp Vault](./17/prd.md)
- [Clerk Documentation](https://clerk.com/docs)
- [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go)
- [HashiCorp Vault Documentation](https://developer.hashicorp.com/vault/docs)
- [Vault Go Client](https://github.com/hashicorp/vault/tree/main/api)


