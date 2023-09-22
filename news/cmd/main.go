package main

import (
	"APIGateway/news/pkg/api"
	"APIGateway/news/pkg/rss"
	"APIGateway/news/pkg/storage"
	"APIGateway/news/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	configURL = "./config.json"
	conn      = "postgres://postgres:PASSWORD@localhost:5432/news" // write your postgresdb password
	newsAddr  = ":8080"
)

func main() {
	psgr, err := postgres.New(conn)
	if err != nil {
		log.Fatal(err)
	}
	a := api.New(psgr)

	chPost := make(chan []storage.Post)
	chErr := make(chan error)

	// Read RSS
	go func() {
		err := rss.GoNews(configURL, chPost, chErr)
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Add posts to db
	go func() {
		for posts := range chPost {
			if err := a.Db.PostMany(posts); err != nil { // check
				chErr <- err
			}
		}
	}()

	// Errors
	go func() {
		for err := range chErr {
			log.Panicln(err)
		}
	}()

	// server launch
	a.Router().Use(Middle)
	log.Printf("News server is launching on %s", newsAddr)
	err = http.ListenAndServe(newsAddr, a.Router())
	if err != nil {
		log.Fatal("News server not launched. Error: ", err)
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
		log.Printf("<-- client ip: %s, method: %s, url: %s, status code: %d %s, trace id: %s", req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)
	})
}
