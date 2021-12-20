package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	bindAddr     string
	authfilePath string
	zonefilePath string
	hookPath     string
)

func init() {
	cobra.OnInitialize(initConfig)
	command.Flags().StringVarP(&bindAddr, "bind", "b", "localhost:8080", "address and port to bind to")
	command.Flags().StringVarP(&authfilePath, "passwd", "p", ".htpasswd", "path to .htpasswd file")
	command.Flags().StringVarP(&zonefilePath, "zone", "z", "ddns.zone", "path to DNS zonefile")
	command.Flags().StringVarP(&hookPath, "hook", "H", "", "full path to command/script to run after updating zonefile")
}

func initConfig() {
	if os.Getenv("JOURNAL_STREAM") != "" {
		log.Default().SetFlags(0)
	}
}

func main() {
	if err := command.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func alwaysMatch(_ *http.Request, _ *mux.RouteMatch) bool {
	return true
}

func respondWithError(w http.ResponseWriter, title string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s\n", title)
}
