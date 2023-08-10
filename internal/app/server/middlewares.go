package server

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.userStore.User().FindByID(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyReqID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"remote_addr: %s, request_id: %s, started %s %s",
			r.RemoteAddr,
			r.Context().Value(ctxKeyReqID),
			r.Method,
			r.RequestURI)
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"remote_addr: %s, request_id: %s, completed in %v",
			r.RemoteAddr,
			r.Context().Value(ctxKeyReqID),
			time.Now().Sub(start))
	})
}
