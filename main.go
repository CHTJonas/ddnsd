package main

import (
	"fmt"
	"net/http"

	auth "github.com/abbot/go-http-auth"
)

func handle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "<html><body><h1>It works!</h1><p>Username: %s</p></body></html>", r.Username)
}

func main() {
	provider := auth.HtpasswdFileProvider(".htpasswd")
	authenticator := auth.NewBasicAuthenticator("ddnsd", provider)
	http.HandleFunc("/", authenticator.Wrap(handle))
	http.ListenAndServe(":8080", nil)
}
