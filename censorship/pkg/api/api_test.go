package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentHandler(t *testing.T) {
	api := New()
	var testBody = []byte(`{"newsID": 35,"content": "Test"}`)
	var testBody1 = []byte(`{"newsID": 35,"content": "qwerty"}`)
	var testBody2 = []byte(`{"newsID": 35,"content": "asdfgh"}`)
	var testBody3 = []byte(`{"newsID": 35,"content": "zxcvbn"}`)
	req := httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody))
	rr := httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusOK)
	}
	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody1))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusBadRequest)
	}
	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody2))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusBadRequest)
	}
	req = httptest.NewRequest(http.MethodPost, "/comments/add", bytes.NewBuffer(testBody3))
	rr = httptest.NewRecorder()
	api.Router().ServeHTTP(rr, req)
	if !(rr.Code == http.StatusBadRequest) {
		t.Errorf("Wrong code: Got %d, Wanted %d", rr.Code, http.StatusBadRequest)
	}
}
