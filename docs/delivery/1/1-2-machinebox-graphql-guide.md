# 1-2 machinebox/graphql Guide (2025-09-05)

## Install
```bash
go get github.com/machinebox/graphql@v0.3.0
```

## Basic Usage
```go
import (
	"context"
	"net/http"
	"github.com/machinebox/graphql"
)

client := graphql.NewClient("https://backboard.railway.app/graphql/v2", &http.Client{})
req := graphql.NewRequest(`query { viewer { id } }`)
req.Header.Set("Authorization", "Bearer "+token)
var resp struct{ Viewer struct{ ID string } }
if err := client.Run(context.Background(), req, &resp); err != nil { /* handle */ }
```

## Notes
- Set `Authorization` header on each request.
- Wrap calls with retry/backoff logic at a higher level if needed.

## Reference
- `https://github.com/machinebox/graphql`
