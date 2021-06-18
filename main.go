package main

import (
	"fmt"
	"net/http"

	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
)

func main() {
	provider := auth.HtpasswdFileProvider(".htpasswd")
	authenticator := auth.NewBasicAuthenticator("ddnsd", provider)
	r := mux.NewRouter()
	r.Path("/update").Methods("POST").Handler(authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		fmt.Fprintf(w, "<html><body><h1>It works!</h1><p>Username: %s</p></body></html>\n", r.Username)
	}))
	r.Path("/update").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, r, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	})
	r.MatcherFunc(alwaysMatch).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, r, "404 Not Found", http.StatusNotFound)
	})
	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	srv.ListenAndServe()
}

func alwaysMatch(_ *http.Request, _ *mux.RouteMatch) bool {
	return true
}

func respondWithError(w http.ResponseWriter, r *http.Request, title string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s\n", title)
}
