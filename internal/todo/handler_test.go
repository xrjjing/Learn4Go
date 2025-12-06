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

	// 登录获取 token
	loginBody := bytes.NewBufferString(`{"email":"admin@example.com","password":"admin123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/v1/login", loginBody)
	loginReq.Header.Set("Content-Type", "application/json")
	loginRR := httptest.NewRecorder()
	handler.ServeHTTP(loginRR, loginReq)
	if loginRR.Code != http.StatusOK {
		t.Fatalf("login code %d", loginRR.Code)
	}
	var loginResp map[string]any
	_ = json.NewDecoder(loginRR.Body).Decode(&loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("login missing token")
	}
	authHeader := "Bearer " + token

	// create
	body := bytes.NewBufferString(`{"title":"learn go"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/todos", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
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
	req = httptest.NewRequest(http.MethodGet, "/v1/todos", nil)
	req.Header.Set("Authorization", authHeader)
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
	req = httptest.NewRequest(http.MethodPut, "/v1/todos/"+strconv.Itoa(created.ID), reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("toggle code %d", rr.Code)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/v1/todos/"+strconv.Itoa(created.ID), nil)
	req.Header.Set("Authorization", authHeader)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete code %d", rr.Code)
	}
}
