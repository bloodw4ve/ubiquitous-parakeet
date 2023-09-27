package main

import (
	"APIGateway/comments/pkg/api"
	"APIGateway/comments/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	dbURL        = "postgres://postgres:PASSWORD@db_comments:5432/comments" // write your postgresdb password
	commentsAddr = ":8081"
)

func main() {
	psgr, err := postgres.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	a := api.New(psgr)
	a.Router().Use(Middle)
	log.Printf("Commentaries server is launching on %s", commentsAddr)
	err = http.ListenAndServe(commentsAddr, a.Router())
	if err != nil {
		log.Fatal("Commantaries server is not launched. Error:", err)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.Header.Get("X-Request-ID")
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)
		statusCode := lrw.statusCode
		log.Printf("<-- client ip: %s, method %s, url: %s, status code: %d %s, trace id: %s", req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)
	})
}
