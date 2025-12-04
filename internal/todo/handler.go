package todo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Server 封装路由与存储。
type Server struct {
	store *Store
	mux   *http.ServeMux
}

// NewServer 创建带路由的 HTTP 处理器。
func NewServer(store *Store) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

// Handler 返回带日志中间件的处理器。
func (s *Server) Handler() http.Handler {
	return loggingMiddleware(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	s.mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			respondJSON(w, s.store.List(), http.StatusOK)
		case http.MethodPost:
			var body struct {
				Title string `json:"title"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondError(w, http.StatusBadRequest, "invalid json")
				return
			}
			if body.Title == "" {
				respondError(w, http.StatusBadRequest, "title required")
				return
			}
			t := s.store.Create(body.Title)
			respondJSON(w, t, http.StatusCreated)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	s.mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/todos/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		switch r.Method {
		case http.MethodPut:
			var body struct {
				Done bool `json:"done"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondError(w, http.StatusBadRequest, "invalid json")
				return
			}
			t, ok := s.store.Toggle(id, body.Done)
			if !ok {
				respondError(w, http.StatusNotFound, "not found")
				return
			}
			respondJSON(w, t, http.StatusOK)
		case http.MethodDelete:
			if !s.store.Delete(id) {
				respondError(w, http.StatusNotFound, "not found")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})
}

func respondJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, map[string]any{"error": msg}, code)
}

// 简易日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
