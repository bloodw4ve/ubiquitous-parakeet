package api

import (
	"APIGateway/comments/pkg/storage"
	"APIGateway/comments/pkg/storage/postgres"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentHandler(t *testing.T) {
	pg, err := postgres.New("postgres://postgres:PASSWORD@localhost:5432/comments") // write your postgresdb password
	if err != nil {
		t.Fatal(err)
	}
	api := New(pg)
	var testBody = []byte(`{"newsID": 35, "content": "Test"}`)
	req := httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusCreated) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	req = httptest.NewRequest(http.MethodGet, "/comments?news_id=35", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Server response was not decoded: %v", err)
	}
	var data []storage.Comment
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Server response was not decoded: %v", err)
	}
	const wantLen = 1
	if len(data) < wantLen {
		t.Fatalf("Got %d records, Wanted %d", len(data), wantLen)
	}

	// Проверяем неверное обращение к handler-у
	req = httptest.NewRequest(http.MethodPost, "/comments/add", nil)
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)

	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusBadRequest)
	}
}
