package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "ddnsd",
	Short: "Dynamic DNS Daemon",
	Run: func(cmd *cobra.Command, args []string) {
		provider := auth.HtpasswdFileProvider(authfilePath)
		authenticator := auth.NewBasicAuthenticator("ddnsd", provider)
		r := mux.NewRouter()
		r.Path("/ping").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "pong")
		})
		r.Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondWithError(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		})
		r.Path("/update").Methods("POST").Handler(authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
			err := r.ParseForm()
			if err != nil {
				respondWithError(w, "500 Internal Server Error", http.StatusInternalServerError)
				log.Println("Error:", err)
				return
			}
			err = updateResourceRecord(r.Username, r.Form.Get("contents"))
			if err != nil {
				respondWithError(w, "500 Internal Server Error", http.StatusInternalServerError)
				log.Println("Error updating zonefile:", err)
				return
			}
			if hookPath != "" {
				err = callHook(hookPath)
				if err != nil {
					respondWithError(w, "500 Internal Server Error", http.StatusInternalServerError)
					log.Println("Error calling hook:", err)
					return
				}
			}
			respondWithError(w, "202 Accepted", http.StatusAccepted)
		}))
		r.Path("/update").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondWithError(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		})
		r.MatcherFunc(alwaysMatch).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondWithError(w, "404 Not Found", http.StatusNotFound)
		})
		r.Use(serverHeaderMiddleware)
		r.Use(proxyMiddleware)
		r.Use(loggingMiddleware)
		srv := &http.Server{
			Handler: r,
			Addr:    bindAddr,
		}
		log.Println("Starting server...")
		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		signal.Notify(c, syscall.SIGQUIT)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		log.Println("Received shutdown signal!")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Println("Waiting for server to exit...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("Shutdown error:", err.Error())
		}
		log.Println("Bye-bye!")
		os.Exit(0)
	},
}
