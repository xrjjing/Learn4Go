package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// 用户不可操作他人 TODO
func TestRBACOwnership(t *testing.T) {
	s := NewServer(NewStore())
	h := s.Handler()

	// admin 创建一个 todo
	adminToken := loginAndGetToken(t, h, "admin@example.com", "admin123")
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title":"admin todo"}`))
	createReq.Header.Set("Authorization", "Bearer "+adminToken)
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, createReq)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d", w.Code)
	}
	var created Todo
	_ = json.NewDecoder(w.Body).Decode(&created)

	// user 尝试删除 admin 的 todo，应被拒绝
	userToken := loginAndGetToken(t, h, "user@example.com", "user123")
	delReq := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(created.ID), nil)
	delReq.Header.Set("Authorization", "Bearer "+userToken)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, delReq)
	if w2.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w2.Code)
	}
}

func loginAndGetToken(t *testing.T, h http.Handler, email, password string) string {
	body := bytes.NewBufferString(`{"email":"` + email + `","password":"` + password + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login %s failed: %d", email, w.Code)
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	token, _ := resp["token"].(string)
	return token
}
