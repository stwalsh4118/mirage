# Railway GraphQL Operations

This directory contains all GraphQL queries, mutations, and subscriptions for the Railway API.

## Structure

- `queries/` - Read operations (queries)
- `mutations/` - Write operations (mutations)
- `subscriptions/` - Real-time subscriptions (WebSocket)

## Naming Conventions

- Use kebab-case for file names
- Be descriptive: `project-by-id.graphql`, `service-create.graphql`
- Match operation name when possible

## Usage

Files are embedded at compile time using `//go:embed`:

```go
//go:embed queries/queries/project-by-id.graphql
var gqlProjectByID string
```

Then used in functions:

```go
func (c *Client) GetProject(ctx context.Context, id string) (Project, error) {
    gql := gqlProjectByID
    // ... use gql string
}
```

## Adding New Operations

1. Create `.graphql` file in appropriate subdirectory
2. Add `//go:embed` directive in corresponding `.go` file
3. Reference the embedded variable in your function

## Benefits

- Syntax highlighting in IDEs
- Easier to read and maintain
- Can be validated independently
- Better version control diffs
- Serves as API documentation


