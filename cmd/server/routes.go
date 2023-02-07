package main

import (
	"net/http"
)

func (s *Server) Routes() {
	s.Router.HandleFunc("/v1/model", s.HandleModel).Methods(http.MethodPut)
	s.Router.HandleFunc("/v1/validate", s.ValidateAPI).Methods(http.MethodPost)
}
