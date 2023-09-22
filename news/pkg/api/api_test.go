package api

import (
	"APIGateway/news/pkg/storage"
	"APIGateway/news/pkg/storage/postgres"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewsHandler(t *testing.T) {
	db, err := postgres.New("postgres://postgres:PASSWORD@localhost:5432/news") // write your postgresdb password
	if err != nil {
		t.Fatalf("Couldn't connect to PostgresDB: %v", err)
	}
	api := New(db)

	// Get news from db
	req := httptest.NewRequest(http.MethodGet, "/news?=page=2&s=", nil)
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	// Decoding JSON.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Couldn't decode server response: %v", err)
	}
	response := struct {
		Posts      []storage.Post
		Pagination storage.Pagination
	}{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		t.Fatalf("Cpuldn't decode server response: %v", err)
	}

	// Get latest news
	req = httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	// Decoding JSON to Post struct
	b, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Couldn't decode server response: %v", err)
	}
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Couldn't decode server response: %v", err)
	}
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("Got %d records, Wanted >= %d", len(data), wantLen)
	}

	// Get searched news
	req = httptest.NewRequest(http.MethodGet, "/news/search?id=2", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	// Decode JSON to Post struct
	b, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Couldn't decode server response: %v", err)
	}
	var post storage.Post
	err = json.Unmarshal(b, &post)
	if err != nil {
		t.Fatalf("Couldn't decode server response: %v", err)
	}

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/news/qwerty", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusNotFound) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusBadRequest)
	}
}
