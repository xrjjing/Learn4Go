package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xrjjing/Learn4Go/rabbitmq-demo/internal/logstore"
	"github.com/xrjjing/Learn4Go/rabbitmq-demo/internal/rabbit"
)

func newTestMux(t *testing.T) (*http.ServeMux, *logstore.Store, context.CancelFunc) {
	t.Helper()
	cfg := loadConfig()
	cfg.Rabbit.Prefetch = 5

	ctx, cancel := context.WithCancel(context.Background())
	mq := rabbit.NewMock(cfg.Rabbit)
	_ = mq.DeclareTopology(ctx)

	mockMsgs := []rabbit.DemoMessage{
		{ID: "m1", Type: "order.created", Payload: json.RawMessage(`{"orderId":"T1"}`)},
		{ID: "m2", Type: "order.fail", Payload: json.RawMessage(`{"orderId":"T2"}`)},
		{ID: "m3", Type: "order.closed", Payload: json.RawMessage(`{"orderId":"T3"}`), TTLMS: 200},
	}
	store := logstore.New(200)
	mux, err := setupHTTP(ctx, cfg, mq, mockMsgs, store)
	if err != nil {
		t.Fatalf("setupHTTP err: %v", err)
	}
	t.Cleanup(func() {
		cancel()
		mq.Close()
	})
	return mux, store, cancel
}

func TestSendAndConsume(t *testing.T) {
	mux, store, _ := newTestMux(t)

	body := []byte(`{"type":"order.created","payload":{"orderId":"X1"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", rr.Code)
	}
	time.Sleep(50 * time.Millisecond)

	logs := store.List()
	if len(logs) < 2 {
		t.Fatalf("expect at least 2 logs, got %d", len(logs))
	}
	if logs[len(logs)-1].Kind != "consume" {
		t.Fatalf("expect consume log, got %s", logs[len(logs)-1].Kind)
	}
}

func TestFailGoesToDLX(t *testing.T) {
	mux, store, _ := newTestMux(t)

	body := []byte(`{"type":"order.fail","payload":{"orderId":"ERR"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", rr.Code)
	}
	time.Sleep(120 * time.Millisecond)

	found := false
	for _, l := range store.List() {
		if l.Kind == "dlx" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expect dlx log when handler returns error")
	}
}

func TestTTLToDLX(t *testing.T) {
	mux, store, _ := newTestMux(t)

	body := []byte(`{"type":"order.closed","ttl_ms":100,"payload":{"orderId":"TTL"}}`)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expect 200, got %d", rr.Code)
	}
	time.Sleep(250 * time.Millisecond)

	found := false
	for _, l := range store.List() {
		if l.Kind == "dlx" && l.Type == "order.closed" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expect dlx log for ttl message")
	}
}
