### Railway GraphQL API: Operations We Consume

This document lists the exact GraphQL operations and patterns we use with a workspace token. We avoid deprecated `me`/`viewer` entry points and query via root fields.

### Tokens and entry points
- Workspace token required
- Use root `Query` fields; avoid `me`/`viewer`

### Core hierarchy
- Workspace → Project → (Environments, Services, Plugins)

### Queries
- projects(after, before, first, last, includeDeleted, userId, workspaceId)
  - Returns connection: `QueryProjectsConnection { edges { node Project } pageInfo }`
  - Minimal usage:
  ```graphql
  query Projects($first: Int!, $after: String) {
    projects(first: $first, after: $after) {
      edges { node { id name } }
      pageInfo { endCursor hasNextPage }
    }
  }
  ```

- project(id: String!) → Project
  - Minimal usage:
  ```graphql
  query Project($id: ID!) { project(id: $id) { id name } }
  ```

- Expanding project details (relay-style connections on Project):
  - environments(after, before, first, last)
  - services(after, before, first, last)
  - plugins(...) (note: plugin mutations are deprecated on Railway)
  - Example:
  ```graphql
  query ProjectDetails($id: ID!, $first: Int!) {
    project(id: $id) {
      id
      name
      environments(first: $first) { edges { node { id name } } }
      services(first: $first)     { edges { node { id name icon } } }
      plugins(first: $first)      { edges { node { id name } } }
    }
  }
  ```

### Mutations

#### Projects
- projectCreate(input: ProjectCreateInput) → Project
  - Fields include: name, description, defaultEnvironmentName, isPublic, isMonorepo, prDeploys, repo(ProjectCreateRepo), runtime(PublicRuntime), teamId
  - Example:
  ```graphql
  mutation CreateProject($input: ProjectCreateInput!) {
    projectCreate(input: $input) { id name }
  }
  ```

- projectUpdate(id: String!, input: ProjectUpdateInput) → Project
  - Fields include: name, description, isPublic, prDeploys, baseEnvironmentId, botPrEnvironments
  - Example:
  ```graphql
  mutation UpdateProject($id: String!, $input: ProjectUpdateInput!) {
    projectUpdate(id: $id, input: $input) { id name }
  }
  ```

- projectDelete(id: String!) → Boolean
  ```graphql
  mutation DeleteProject($id: String!) { projectDelete(id: $id) }
  ```

#### Environments
- environmentCreate(input: EnvironmentCreateInput) → Environment
  - Required: name, projectId
  - Optional: ephemeral, skipInitialDeploys, stageInitialChanges, applyChangesInBackground, sourceEnvironmentId
  - Example:
  ```graphql
  mutation CreateEnvironment($input: EnvironmentCreateInput!) {
    environmentCreate(input: $input) { id name }
  }
  ```

- environmentDelete(id: String!) → Boolean
  ```graphql
  mutation DeleteEnvironment($id: String!) { environmentDelete(id: $id) }
  ```

- environmentRename(id: String!, input: EnvironmentRenameInput) → Environment
  ```graphql
  mutation RenameEnvironment($id: String!, $input: EnvironmentRenameInput!) {
    environmentRename(id: $id, input: $input) { id name }
  }
  ```

- Advanced (optional): environmentPatch / environmentPatchCommit for staged changes.

#### Services
- serviceCreate(input: ServiceCreateInput) → Service
  - Required: projectId
  - Optional: name, source(ServiceSourceInput), branch, icon, environmentId (scoped create), templateServiceId, registryCredentials, variables(EnvironmentVariables)
  - Example:
  ```graphql
  mutation CreateService($input: ServiceCreateInput!) {
    serviceCreate(input: $input) { id name }
  }
  ```

- serviceUpdate(id: String!, input: ServiceUpdateInput) → Service
  ```graphql
  mutation UpdateService($id: String!, $input: ServiceUpdateInput!) {
    serviceUpdate(id: $id, input: $input) { id name icon }
  }
  ```

- serviceDelete(id: String!, environmentId: String) → Boolean
  ```graphql
  mutation DeleteService($id: String!, $environmentId: String) {
    serviceDelete(id: $id, environmentId: $environmentId)
  }
  ```

#### Plugins (deprecated by Railway)
- pluginCreate / pluginUpdate / pluginDelete exist but are marked deprecated. Prefer database templates/integrations. Only use if legacy support is required.

### Pagination pattern
- All list fields are relay-style: use `after` + `first`, read `pageInfo { endCursor hasNextPage }` and loop until complete.
  ```graphql
  query Projects($first: Int!, $after: String) {
    projects(first: $first, after: $after) {
      edges { node { id name } }
      pageInfo { endCursor hasNextPage }
    }
  }
  ```

### Notes
### Input type references (from Railway introspection)

- ProjectCreateInput
  - defaultEnvironmentName: String
  - description: String
  - isMonorepo: Boolean
  - isPublic: Boolean
  - name: String
  - prDeploys: Boolean
  - repo: ProjectCreateRepo
  - runtime: PublicRuntime
  - teamId: String

- ProjectCreateRepo
  - branch: String!
  - fullRepoName: String!

- ProjectUpdateInput
  - baseEnvironmentId: String
  - botPrEnvironments: Boolean
  - description: String
  - isPublic: Boolean
  - name: String
  - prDeploys: Boolean

- EnvironmentCreateInput
  - applyChangesInBackground: Boolean
  - ephemeral: Boolean
  - name: String!
  - projectId: String!
  - skipInitialDeploys: Boolean
  - sourceEnvironmentId: String
  - stageInitialChanges: Boolean

- EnvironmentRenameInput
  - name: String!

- ServiceCreateInput
  - branch: String
  - environmentId: String
  - icon: String
  - name: String
  - projectId: String!
  - registryCredentials: RegistryCredentialsInput
  - source: ServiceSourceInput
  - templateServiceId: String
  - variables: EnvironmentVariables

- ServiceUpdateInput
  - icon: String
  - name: String

Notes:
- `String!` indicates required.
- Some nested inputs (e.g., ServiceSourceInput, RegistryCredentialsInput) are available; we’ll include them when we integrate those features.

- Use workspace tokens; do not rely on `me`/`viewer`.
- For details pages, request only the fields needed to keep payloads small.


