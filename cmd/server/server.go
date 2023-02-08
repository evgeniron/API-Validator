package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/evgeniron/API-Validator/model"
	"github.com/evgeniron/API-Validator/store"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	db     *store.Store
}

func NewServer(db *store.Store) *Server {
	s := &Server{}
	s.db = db
	s.Routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ValidateAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) HandleModel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models, err := model.Decode(r)
		if err != nil {
			respond(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = model.StoreModels(s.db, models)
		if err != nil {
			respond(w, r, http.StatusBadRequest, err.Error())
			return
		}
		respond(w, r, http.StatusOK, nil)
	}
}

type Response struct {
	valid bool
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		log.Println("respond:", &Response{true})
	}
}
