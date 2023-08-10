package server

import (
	"encoding/json"
	"errors"
	"github.com/GlebSr/nuhaiShort/internal/app/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *server) handleUrlRedirection() http.HandlerFunc {
	type redirect struct {
		Location string `json:"Location"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url, err := s.urlStore.Url().FindByID(vars["key"])
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		err = s.urlStore.Url().RegisterRedirect(url)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		w.Header().Set("Location", url.LongUrl)
		s.respond(w, r, http.StatusTemporaryRedirect, nil)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.userStore.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.userStore.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		session, err := s.sessionStore.New(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		err = s.sessionStore.Save(r, w, session)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleNewUrl() http.HandlerFunc {
	type request struct {
		LongUrl string `json:"long_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		url := &model.Url{
			LongUrl:         req.LongUrl,
			UserId:          r.Context().Value(ctxKeyUser).(*model.User).ID,
			RedirectsNumber: 0,
		}
		for {

			if err := url.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			err := url.MakeShort()
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, errors.New("generating error"))
				return
			}
			select {
			case <-r.Context().Done():
				s.error(w, r, http.StatusGatewayTimeout, errors.New("server time out"))
				return
			default:
				err := s.urlStore.Url().Create(url)
				if err != nil {
					continue
				}
				s.respond(w, r, http.StatusOK, url)
				return
			}
		}

	}
}

func (s *server) handleNewCoolUrl() http.HandlerFunc {
	type request struct {
		Key      string `json:"key"`
		LongUrl  string `json:"long_url"`
		ShortUrl string `json:"short_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Key != "key" {
			s.error(w, r, http.StatusMethodNotAllowed, errors.New("i neeeeeed key"))
			return
		}
		url := &model.Url{
			LongUrl:         req.LongUrl,
			UserId:          r.Context().Value(ctxKeyUser).(*model.User).ID,
			RedirectsNumber: 0,
			ID:              req.ShortUrl,
		}
		err := s.urlStore.Url().Create(url)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, url)
	}
}

func (s *server) handleDeleteUrl() http.HandlerFunc {
	type request struct {
		Url string `json:"short_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		url, err := s.urlStore.Url().FindByID(req.Url)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if url.UserId != r.Context().Value(ctxKeyUser).(*model.User).ID {

			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err = s.urlStore.Url().Delete(req.Url)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleOneUrlInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url, err := s.urlStore.Url().FindByID(vars["key"])
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if url.UserId != r.Context().Value(ctxKeyUser).(*model.User).ID {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, url)
	}
}

func (s *server) handleUrlsInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := s.urlStore.Url().FindByUserId(r.Context().Value(ctxKeyUser).(*model.User).ID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, urls)
	}
}
