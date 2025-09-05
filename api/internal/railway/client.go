package railway

import (
	"context"
	"net/http"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/machinebox/graphql"
)

const (
	DefaultEndpoint = "https://backboard.railway.app/graphql/v2"
)

// Client wraps GraphQL calls to Railway with retries and auth.
type Client struct {
	endpoint string
	token    string
	httpc    *http.Client
}

func NewClient(endpoint, token string, httpc *http.Client) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if httpc == nil {
		httpc = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{endpoint: endpoint, token: token, httpc: httpc}
}

// execute runs a GraphQL operation with retries on errors within a bounded window.
func (c *Client) execute(ctx context.Context, gql string, vars map[string]any, out any) error {
	client := graphql.NewClient(c.endpoint, graphql.WithHTTPClient(c.httpc))
	req := graphql.NewRequest(gql)
	for k, v := range vars {
		req.Var(k, v)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	operation := func() error {
		return client.Run(ctx, req, out)
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 300 * time.Millisecond
	bo.MaxInterval = 2 * time.Second
	bo.MaxElapsedTime = 10 * time.Second

	if err := backoff.Retry(func() error { return operation() }, backoff.WithContext(bo, ctx)); err != nil {
		return err
	}
	return nil
}
