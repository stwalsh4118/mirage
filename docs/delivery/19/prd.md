# PBI-19: Import and Manage External Railway Environments

## Overview
Enable users to import Railway environments that existed before connecting to Mirage, allowing full management capabilities through Mirage's interface. This closes the gap between Mirage-created resources and pre-existing Railway infrastructure.

## Problem Statement
Currently, Mirage only tracks and allows operations on resources it created (stored in the database with user associations). When users connect their Railway account to Mirage:
- Pre-existing Railway resources are visible through the browsing interface (PBI 9)
- Users cannot perform any management operations on these external resources
- This creates a fragmented experience requiring context-switching between Mirage and Railway UIs
- Users cannot consolidate their Railway infrastructure management under Mirage

## User Stories
- As a user, I want to see which Railway environments are managed by Mirage and which are external, so I can understand what I can control.
- As a user, I want to import an external Railway environment into Mirage, so I can manage it through Mirage's interface.
- As a user, I want to import multiple environments at once, so I don't have to import them one by one.
- As a user, I want to be warned if importing an environment will cause conflicts, so I can make informed decisions.
- As a user, I want all environment metadata (variables, services, configurations) preserved during import, so nothing is lost.

## Technical Approach

### Visual Distinction
- Display a badge/indicator on environments showing their management status:
  - "Managed by Mirage" - Created by or imported into Mirage
  - "External" - Exists in Railway but not tracked in Mirage DB
- Use consistent visual language across project views and environment lists

### Import Flow (Environment-Level)
1. **Single Import**:
   - Add "Import" button on external environments in the UI
   - On click, validate for conflicts (see validation section)
   - If valid, fetch full environment metadata from Railway API
   - Store in Mirage database with proper user association
   - Update UI to reflect new managed status

2. **Mass Import**:
   - Add "Import All" or "Import Selected" capability on project view
   - Allow multi-select of external environments
   - Batch validation check for all selected environments
   - Display validation results (which will succeed/fail)
   - Proceed with import for validated environments
   - Show progress indicator for batch operation

### Validation Rules
Before importing an environment, validate:
1. **Uniqueness**: No existing Mirage environment with the same Railway environment ID
2. **Service Conflicts**: Check if any services within the environment have Railway IDs that conflict with existing Mirage-tracked services
3. **Metadata Integrity**: Verify that all required Railway metadata is accessible (environment variables, service configurations)
4. **User Authorization**: Confirm user has valid Railway API token with sufficient permissions
5. **Project Association**: Ensure the environment's parent project is accessible

If validation fails:
- Display clear error messages explaining the conflict
- Suggest remediation steps (e.g., "Environment already managed", "Insufficient API permissions")
- Do not proceed with import

### Metadata Capture
During import, capture and store:
- Environment Railway ID and name
- Associated project Railway ID
- All services within the environment:
  - Service Railway IDs
  - Service names and source configurations
  - Build settings (if available)
  - Deployment configurations
- Environment variables (names and encrypted values via Vault)
- Environment creation date and metadata
- User association (linking to current Mirage user)
- Import timestamp

### Post-Import Operations
Once imported:
- All standard Mirage CRUD operations become available
- Environment appears in user's managed environments list
- Can modify environment variables, add/remove services
- Can delete environment (with confirmation)
- Audit logs track all operations on imported environments

### Database Schema Considerations
Ensure existing schema supports:
- Railway environment ID as unique identifier
- Flag or status field to distinguish import vs. created origins
- Proper foreign key relationships for user ownership
- Service-to-environment relationships

## UX/UI Considerations

### Visual Design
- **Badge System**: 
  - "Managed" badge: Green/teal color scheme
  - "External" badge: Gray/neutral color scheme
- **Import Button**: 
  - Prominent but not primary action
  - Located near environment name/header
  - Disabled state if validation fails
- **Mass Import Interface**:
  - Checkbox selection for external environments
  - Sticky action bar with "Import Selected" button
  - Clear count of selected environments

### User Feedback
- Loading states during import process
- Success notifications with confirmation
- Detailed error messages with actionable guidance
- Progress indicators for batch imports
- Confirmation dialogs for mass import operations

### Information Architecture
- Clear separation between managed and external resources in lists
- Filter/sort options to show only managed or external environments
- Breadcrumb context showing import status in environment details

## Acceptance Criteria
1. ✅ Visual distinction between Mirage-managed and external Railway environments is clear and consistent across all views
2. ✅ Users can import a single external environment with one click
3. ✅ Users can select and import multiple external environments in a single operation
4. ✅ Validation checks run before import and prevent conflicts
5. ✅ Validation errors display clear, actionable messages to users
6. ✅ All environment metadata (services, variables, configs) is captured during import
7. ✅ Imported environments support full CRUD operations identical to Mirage-created environments
8. ✅ Import operations are associated with the current user's account
9. ✅ Database properly stores Railway resource IDs and maintains referential integrity
10. ✅ Audit trail captures all import operations with timestamps and user attribution

## Dependencies
- PBI 9: Railway project browsing (provides visibility into external resources)
- PBI 7: Persistence and RBAC (database schema for resource storage)
- PBI 17: HashiCorp Vault (secure storage of imported environment variables)
- Railway API: Full read access to environment and service metadata
- Database schema supporting Railway ID storage and user associations

## Open Questions
1. Should we support "un-importing" (removing from Mirage management but keeping in Railway)?
2. How do we handle environments that are modified in Railway directly after import (sync/conflict detection)?
3. Should there be a role/permission requirement for import operations?
4. Do we need to track import history (who imported, when) separately from general audit logs?
5. Should we provide a "preview" mode showing what would be imported before committing?

## Related Tasks
[View in Backlog](../backlog.md#19)

