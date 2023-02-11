package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/evgeniron/API-Validator/model"
	"github.com/evgeniron/API-Validator/store"
	"github.com/evgeniron/API-Validator/validate"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	db     *store.Store
}

func NewServer(db *store.Store) *Server {
	s := &Server{}
	s.db = db
	s.router = mux.NewRouter()
	s.Routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) ValidateEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parse request
		endpoint, err := validate.ParseEndpoint(r)
		if err != nil {
			respond(w, r, http.StatusBadRequest, err.Error())
		}

		// Load model
		model, err := model.GetModel(s.db, endpoint.Path, endpoint.Method)
		if err != nil {
			var e *store.RecordNotFoundError
			if errors.As(err, &e) {
				// we can add metrics/logs here for records not found
				respond(w, r, http.StatusOK, nil)
				return
			}
			respond(w, r, http.StatusInternalServerError, nil)
		}

		// Validate endpoint
		report, err := validate.ValidateReport(endpoint, model)
		if err != nil {
			// we can add metrics/logs/traces here for errors validating record
			respond(w, r, http.StatusOK, nil)
			return
		}

		// Write response
		respond(w, r, http.StatusOK, report)
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

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
