package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main()  {
	var (
		addr = flag.String("addr", ":8080", "endpoint address")
	)
	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CONFIG"))
	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	defer client.Close()

	s := &Server{db: client}
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
	log.Println("Starting webserver on", *addr)
	http.ListenAndServe(":8080", mux)
	log.Println("Stopping...")
}

// Server is the API server
type Server struct {
	db *firestore.Client
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
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