package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

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
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		httpInfo := fmt.Sprintf("\"%s %s %s\"", r.Method, r.URL.Path, r.Proto)
		refererInfo := fmt.Sprintf("\"%s\"", r.Referer())
		if refererInfo == "\"\"" {
			refererInfo = "\"-\""
		}
		uaInfo := fmt.Sprintf("\"%s\"", r.UserAgent())
		if uaInfo == "\"\"" {
			uaInfo = "\"-\""
		}
		log.Println(r.RemoteAddr, httpInfo, lrw.statusCode, refererInfo, uaInfo)
	})
}
