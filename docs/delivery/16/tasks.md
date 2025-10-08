# Tasks for PBI 16: Clerk Authentication & Resource Ownership (Core)

This document lists all tasks associated with PBI 16.

**Parent PBI**: [PBI 16: Clerk Authentication & Resource Ownership (Core)](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--------------------------------------- | :------- | :--------------------------------- |
| 16-1 | [Research and document Clerk Go SDK](./16-1.md) | Proposed | Research Clerk Go SDK v2 API, JWT verification patterns, and create implementation guide |
| 16-2 | [Setup Clerk application and configure webhooks](./16-2.md) | Proposed | Create Clerk app, configure dev environment, set up webhook endpoints |
| 16-3 | [Create User model and database migration](./16-3.md) | Proposed | Define User GORM model, create migration, add indexes |
| 16-4 | [Add UserID foreign keys to existing models and migration](./16-4.md) | Proposed | Add UserID to Environment, Service, EnvironmentMetadata models, create migration, delete test data |
| 16-5 | [Implement Clerk JWT verification middleware](./16-5.md) | Proposed | Create RequireAuth middleware with JWT verification, user loading, and context helpers |
| 16-6 | [Implement Clerk webhook handler](./16-6.md) | Proposed | Create webhook endpoint with signature verification and event handlers (user.created, user.updated, user.deleted) |
| 16-7 | [Update Environment controller with ownership checks](./16-7.md) | Proposed | Add UserID to environment creation, filter queries by user, implement 403 checks |
| 16-8 | [Update Services controller with ownership checks](./16-8.md) | Proposed | Add UserID to service creation, filter queries by user, verify environment ownership |
| 16-9 | [Implement user profile API endpoints](./16-9.md) | Proposed | Create GET/PATCH /api/v1/users/me and resource list endpoints |
| 16-10 | [Apply authentication middleware to all protected routes](./16-10.md) | Proposed | Apply RequireAuth to all routes except healthz and webhooks, update server.go |
| 16-11 | [Install and configure Clerk frontend SDK](./16-11.md) | Done | Install @clerk/nextjs, configure ClerkProvider, set up environment variables |
| 16-12 | [Create sign-in and sign-up pages](./16-12.md) | Done | Implement /sign-in and /sign-up pages with Clerk components |
| 16-13 | [Add authentication state management and JWT interceptor](./16-13.md) | Done | Update API client to send JWT tokens, handle 401 responses, add auth state to UI |
| 16-14 | [Implement user profile UI components](./16-14.md) | Proposed | Add user profile dropdown to header with avatar, name, email, sign-out |
| 16-15 | [Add protected route guards to frontend](./16-15.md) | Proposed | Implement auth checks for protected pages, redirect unauthenticated users |
| 16-16 | [Update environment and service creation flows with user context](./16-16.md) | Proposed | Ensure wizard and API calls work with authenticated user, remove UserID from frontend payloads |
| 16-17 | [Write unit tests for authentication middleware and helpers](./16-17.md) | Proposed | Test JWT verification, user loading, context helpers, error cases |
| 16-18 | [Write integration tests for webhook handlers](./16-18.md) | Proposed | Test user CRUD via webhooks, signature verification, idempotency |
| 16-19 | [Write integration tests for ownership checks](./16-19.md) | Proposed | Test resource isolation, 403 responses, cross-user access prevention |
| 16-20 | [E2E CoS Test: Complete authentication and ownership flow](./16-20.md) | Proposed | End-to-end test of sign-up, resource creation, ownership isolation, profile management |

