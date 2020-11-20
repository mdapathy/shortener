package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mdapathy/url-shortener/url/controller"
	"net/http"
)

type ApiServer struct {
	Port       int
	Controller controller.Controller
	server     *http.Server
}

func (s *ApiServer) Create() http.Handler {
	if s.Port == 0 {
		panic("port for the server was not defined")
	}
	r := mux.NewRouter()

	r.HandleFunc("/url/{key}", s.Controller.HandleUrlRedirect).
		Methods("GET")

	r.HandleFunc("/url/{key}", s.Controller.HandleUriDelete).
		Methods("DELETE")

	r.HandleFunc("/shorten", s.Controller.HandleUriCreate).
		Methods("POST")

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: r,
	}

	return r
}

func (s *ApiServer) Start() error {
	s.Create()
	return s.server.ListenAndServe()
}

func (s *ApiServer) Stop() error {
	if s.server == nil {
		return fmt.Errorf("no server found")
	}
	return s.server.Shutdown(context.Background())
}
