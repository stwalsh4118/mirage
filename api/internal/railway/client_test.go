package railway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestExecute_RetriesOnTransientErrors(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		if c < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors":[{"message":"boom"}]}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))

	defer ts.Close()

	c := NewClient(ts.URL, "", ts.Client())
	var out any
	err := c.execute(context.Background(), `query { ok }`, nil, &out)
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", calls)
	}
}
