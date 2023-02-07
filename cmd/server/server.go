package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func newServer() *Server {
	s := &Server{}
	s.Routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) ValidateAPI(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleModel(w http.ResponseWriter, r *http.Request) {

}
