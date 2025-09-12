# [10-5] Railway API Guide (2025-09-11)

[Back to task list](./tasks.md)

This guide aggregates authoritative references and safe usage notes for interacting with Railway to create Projects, Environments, and Services. It avoids guessing mutation names and focuses on verified approaches and schema introspection.

## References
- Railway Environments (concepts): `https://docs.railway.com/reference/environments`
- Railway Environments (develop workflows): `https://docs.railway.com/develop/environments`
- Railway Projects (concepts/UX): `https://docs.railway.com/develop/projects`

Note: Railway’s public documentation primarily covers concepts and UI. Programmatic API details (GraphQL operations) are not fully documented publicly. We will confirm concrete operations via schema introspection on the GraphQL endpoint available to authenticated users, or fall back to officially supported SDKs/endpoints if exposed.

## Authentication
- Token: Personal or service token set as `RAILWAY_API_TOKEN` on the backend.
- Header: `Authorization: Bearer <token>`.

## API Access Strategy
1. Identify the GraphQL endpoint used by the Railway app (historically `https://backboard.railway.app/graphql` or equivalent). Confirm via network inspection and update this guide with the exact host during implementation.
2. Perform GraphQL introspection to list available mutations and inputs. Store key snippets below once verified.
3. Prefer higher-level operations if available (e.g., create project, create environment). If only lower-level primitives exist, compose them.
4. Ensure idempotency: pass client request IDs where supported or implement client-side deduplication and retries with backoff.

## Confirmed Operations (per schema examples provided)

Note: Field names and input object shapes below reflect working mutations provided during design. Preserve nullability as shown and supply concrete values in implementation.

### Create Project
Mutation name: `projectCreate`

Input example:
```graphql
mutation ProjectCreate($defaultEnvironmentName: String) {
  projectCreate(input: { defaultEnvironmentName: $defaultEnvironmentName }) {
    baseEnvironmentId
    botPrEnvironments
    createdAt
    deletedAt
    description
    expiredAt
    id
    isPublic
    isTempProject
    name
    prDeploys
    subscriptionPlanLimit
    subscriptionType
    teamId
    updatedAt
    workspaceId
  }
}
```

Variables example:
```json
{ "defaultEnvironmentName": "production" }
```

### Create Environment
Mutation name: `environmentCreate`

Input example:
```graphql
mutation EnvironmentCreate($name: String, $projectId: ID!) {
  environmentCreate(input: { name: $name, projectId: $projectId }) {
    createdAt
    deletedAt
    id
    isEphemeral
    name
    projectId
    unmergedChangesCount
    updatedAt
  }
}
```

Variables example:
```json
{ "projectId": "<PROJECT_ID>", "name": "staging" }
```

### Create Service
Mutation name: `serviceCreate`

Input example:
```graphql
mutation ServiceCreate($name: String!, $environmentId: ID!, $branch: String, $repo: String) {
  serviceCreate(
    input: { branch: $branch, environmentId: $environmentId, source: { repo: $repo }, name: $name }
  ) {
    createdAt
    deletedAt
    featureFlags
    icon
    id
    name
    projectId
    templateServiceId
    templateThreadSlug
    updatedAt
  }
}
```

Variables example:
```json
{ "name": "api", "environmentId": "<ENV_ID>", "branch": "main", "repo": "github.com/org/repo" }
```

## Example Flows

Create a new project (with default environment name) and then an environment:
```graphql
mutation ProjectCreate($defaultEnvironmentName: String) {
  projectCreate(input: { defaultEnvironmentName: $defaultEnvironmentName }) { id name baseEnvironmentId }
}

mutation EnvironmentCreate($projectId: ID!, $name: String) {
  environmentCreate(input: { projectId: $projectId, name: $name }) { id name projectId }
}
```

Add an environment to an existing project:
```graphql
mutation EnvironmentCreate($projectId: ID!, $name: String) {
  environmentCreate(input: { projectId: $projectId, name: $name }) { id name projectId }
}
```

Provision a service from a repo/branch:
```graphql
mutation ServiceCreate($environmentId: ID!, $name: String!, $repo: String, $branch: String) {
  serviceCreate(input: { environmentId: $environmentId, name: $name, source: { repo: $repo }, branch: $branch }) {
    id
    name
    projectId
  }
}
```

Set environment variables:
```graphql
mutation SetVars($projectId: ID!, $environmentId: ID!, $vars: [KeyValueInput!]!) {
  setEnvironmentVariables(input: { scope: { projectId: $projectId, environmentId: $environmentId }, variables: $vars })
}
```

## Verbatim mutations (as provided)

These are the exact snippets provided during design, preserved verbatim for reference.

```graphql
mutation EnvironmentCreate {
    environmentCreate(input: { name: null, projectId: null }) {
        createdAt
        deletedAt
        id
        isEphemeral
        name
        projectId
        unmergedChangesCount
        updatedAt
    }
}
```

```graphql
mutation ProjectCreate {
    projectCreate(
        input: { defaultEnvironmentName: null }
    ) {
        baseEnvironmentId
        botPrEnvironments
        createdAt
        deletedAt
        description
        expiredAt
        id
        isPublic
        isTempProject
        name
        prDeploys
        subscriptionPlanLimit
        subscriptionType
        teamId
        updatedAt
        workspaceId
    }
}
```

```graphql
    serviceCreate(
        input: { branch: null, environmentId: null, source: { repo: null }, name: null }
    ) {
        createdAt
        deletedAt
        featureFlags
        icon
        id
        name
        projectId
        templateServiceId
        templateThreadSlug
        updatedAt
    }
```

## Error Handling and Retries
- Classify: validation vs transient vs permission.
- Retries: exponential backoff on 5xx/network; do not retry 4xx validation errors.
- Partial failures: surface resource IDs created so far; offer cleanup or resume.

## TODO (to be completed during implementation)
- Replace pseudo operations with actual mutation names, inputs, and endpoint URL after schema introspection.
- Add cURL examples and real request/response fragments captured from a dev token.
- Document pagination/list queries to populate the Existing Project selector.
