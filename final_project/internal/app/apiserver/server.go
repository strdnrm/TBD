package apiserver

import (
	"context"
	"encoding/json"
	"final_project/internal/app/model"
	"final_project/internal/app/store"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type server struct {
	router *mux.Router
	logger *zap.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	s := &server{
		router: mux.NewRouter(),
		logger: l,
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.hanldeUsersCreate()).Methods("POST")
}

func (s *server) hanldeUsersCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Login:    req.Login,
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
			Surname:  req.Surname,
		}

		if err := s.store.User().Create(context.Background(), u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
