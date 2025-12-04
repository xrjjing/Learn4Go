package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestTodoFlow(t *testing.T) {
	s := NewServer(NewStore())
	handler := s.Handler()

	// create
	body := bytes.NewBufferString(`{"title":"learn go"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201 got %d", rr.Code)
	}
	var created Todo
	_ = json.NewDecoder(rr.Body).Decode(&created)
	if created.Title != "learn go" {
		t.Fatalf("title mismatch")
	}

	// list
	req = httptest.NewRequest(http.MethodGet, "/todos", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list code %d", rr.Code)
	}
	var list []Todo
	_ = json.NewDecoder(rr.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("want 1 item, got %d", len(list))
	}

	// toggle
	reqBody := bytes.NewBufferString(`{"done":true}`)
	req = httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(created.ID), reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("toggle code %d", rr.Code)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(created.ID), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete code %d", rr.Code)
	}
}
