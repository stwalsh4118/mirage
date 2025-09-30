# Railway Public API v2 – Quick Guide (2025-09-10)

[Source: Railway Public API Guide](https://docs.railway.com/guides/public-api) • Endpoint: `https://backboard.railway.com/graphql/v2`

## Auth
- Header: `Authorization: Bearer <RAILWAY_API_TOKEN>`

## Test Query
```bash
curl --request POST \
  --url https://backboard.railway.com/graphql/v2 \
  --header "Authorization: Bearer $RAILWAY_API_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{"query":"query { viewer { id } }"}'
```

## Common Operations (subject to schema)
- Create environment (project-scoped): mutation likely named similar to `createEnvironment` taking `projectId` and `name`, returning `environment { id }`.
- Delete environment: mutation similar to `deleteEnvironment(id: ID!) { id }`.
- Get environment status: query `environment(id: ID!) { id status }` or deployment state under project/service.

## Notes
- Use endpoint override if needed: env `RAILWAY_GRAPHQL_ENDPOINT`.
- Normalize remote states to our UI set: creating | active | destroying | error | unknown.

## References
- Railway Public API: https://docs.railway.com/guides/public-api
- Endpoint: https://backboard.railway.com/graphql/v2







