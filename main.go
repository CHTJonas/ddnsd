package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
	zonefile "github.com/bwesterb/go-zonefile"
	"github.com/gorilla/mux"
)

const (
	authfilePath = ".htpasswd"
	zonefilePath = "test.zone"
)

func main() {
	provider := auth.HtpasswdFileProvider(authfilePath)
	authenticator := auth.NewBasicAuthenticator("ddnsd", provider)
	r := mux.NewRouter()
	r.Path("/update").Methods("POST").Handler(authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		err := r.ParseForm()
		if err != nil {
			respondWithError(w, "500 Internal Server Error", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		err = updateResourceRecord(r.Username, r.Form.Get("contents"))
		if err != nil {
			respondWithError(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		respondWithError(w, "202 Accepted", http.StatusAccepted)
	}))
	r.Path("/update").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	})
	r.MatcherFunc(alwaysMatch).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, "404 Not Found", http.StatusNotFound)
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

func respondWithError(w http.ResponseWriter, title string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s\n", title)
}

func updateResourceRecord(username, contents string) error {
	// Load zonefile
	data, err := ioutil.ReadFile(zonefilePath)
	if err != nil {
		return err
	}

	// Parse zonefile
	zf, err := zonefile.Load(data)
	if err != nil {
		return err
	}

	// Update RR
	for _, e := range zf.Entries() {
		if !bytes.Equal(e.Domain(), []byte(username)) {
			continue
		}
		if !bytes.Equal(e.Type(), []byte("TXT")) {
			return errors.New("resource record type in zonefile is not TXT")
		}
		e.SetValue(0, []byte(contents))
		fh, err := os.OpenFile(zonefilePath, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		_, err = fh.Write(zf.Save())
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("could not find resource record in zonefile")
}
