package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func serverHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "https://github.com/CHTJonas/ddnsd")
		next.ServeHTTP(w, r)
	})
}

func proxyMiddleware(next http.Handler) http.Handler {
	return handlers.ProxyHeaders(next)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpAction := fmt.Sprintf("\"%s %s %s\"", r.Method, r.URL.Path, r.Proto)
		fmt.Println(r.RemoteAddr, httpAction)
		next.ServeHTTP(w, r)
	})
}
