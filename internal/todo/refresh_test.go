package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 覆盖 refresh 正常与重复使用的错误路径
func TestRefreshToken(t *testing.T) {
	s := NewServer(NewStore())
	handler := s.Handler()

	// login 获取 refresh
	loginBody := bytes.NewBufferString(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", loginBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("login status %d", rr.Code)
	}
	var loginResp map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&loginResp)
	refresh, ok := loginResp["refresh_token"].(string)
	if !ok || refresh == "" {
		t.Fatalf("missing refresh token")
	}

	refreshReqBody := bytes.NewBufferString(`{"refresh_token":"` + refresh + `"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/refresh", refreshReqBody)
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("refresh status %d", rr2.Code)
	}

	// 再次使用旧 refresh 应该失败（被旋转作废）
	req3 := httptest.NewRequest(http.MethodPost, "/refresh", refreshReqBody)
	req3.Header.Set("Content-Type", "application/json")
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	if rr3.Code == http.StatusOK {
		t.Fatalf("expected refresh reuse to fail")
	}
}
