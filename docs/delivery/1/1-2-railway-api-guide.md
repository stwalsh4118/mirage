# 1-2 Railway GraphQL API Guide (2025-09-05)

## Endpoint
- URL: `https://backboard.railway.app/graphql/v2`
- Auth: `Authorization: Bearer <RAILWAY_API_TOKEN>`
- Content-Type: `application/json`

## Notes
- Use POST requests with GraphQL `query` or `mutation` payloads.
- For reliability, implement retries on 429 and 5xx responses with exponential backoff and jitter.

## Example Request
```http
POST /graphql/v2 HTTP/1.1
Host: backboard.railway.app
Authorization: Bearer <token>
Content-Type: application/json

{"query":"query { viewer { id } }"}
```

## Mutations (to confirm in implementation)
- Create Environment: mutation name and input shape to be confirmed during implementation.
- Destroy Environment: mutation name and input shape to be confirmed during implementation.

## References
- Railway docs: `https://docs.railway.app` (public docs)
- Known endpoint usage in community examples references `backboard.railway.app/graphql/v2`
