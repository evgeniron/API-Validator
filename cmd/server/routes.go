package main

import (
	"net/http"
)

const (
	modelRoute      = "/v1/model"
	validationRoute = "/v1/validate"
)

func (s *Server) Routes() {
	s.router.HandleFunc(modelRoute, s.HandleModel()).Methods(http.MethodPut)
	s.router.HandleFunc(validationRoute, s.ValidateEndpoint()).Methods(http.MethodPost)
}
