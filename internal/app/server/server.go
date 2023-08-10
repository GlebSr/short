package server

import (
	"encoding/json"
	"github.com/GlebSr/nuhaiShort/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

type server struct {
	router       *mux.Router
	userStore    store.UserStore
	urlStore     store.UrlStore
	sessionStore sessions.Store
}

type ctxKey int8

const (
	sessionName        = "dora"
	ctxKeyUser  ctxKey = iota
	ctxKeyReqID
)

func newServer(userStore store.UserStore, urlStore store.UrlStore, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		userStore:    userStore,
		urlStore:     urlStore,
		sessionStore: sessionStore,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	v1Router := s.router.PathPrefix("/v1").Subrouter()
	v1Router.Use(s.logRequest)
	v1Router.HandleFunc("/registration", s.handleUsersCreate()).Methods("POST")
	v1Router.HandleFunc("/login", s.handleSessionCreate()).Methods("POST")
	private := v1Router.PathPrefix("/account").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/new-url", s.handleNewUrl()).Methods("POST")
	private.HandleFunc("/new-special-url", s.handleNewCoolUrl()).Methods("POST")
	private.HandleFunc("/delete-url", s.handleDeleteUrl()).Methods("DELETE")
	private.HandleFunc("/urls-information", s.handleUrlsInformation()).Methods("GET")
	private.HandleFunc("/urls-information/{key}", s.handleOneUrlInformation()).Methods("GET")
	s.router.HandleFunc("/{key}", s.handleUrlRedirection()).Methods("GET")
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
