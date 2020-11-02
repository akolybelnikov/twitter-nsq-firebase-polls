package main

import (
	"context"
	"net/http"
)

func main()  {
	
}

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

// APIKey extracts and returns the key, given the context
func APIKey(ctx context.Context) (string, bool)  {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc  {
	return func(rw http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(rw, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(rw, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool  {
	return key == "abc123"
}