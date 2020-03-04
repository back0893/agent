package http

import (
	"context"
	"net/http"
)

type Server struct {
	Server *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		Server: &http.Server{Addr: addr, Handler: nil},
	}
}
func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}

func (s *Server) AddHandler(patter string, fn func(writer http.ResponseWriter, request *http.Request)) {
	http.HandleFunc(patter, fn)
}
func (s *Server) Close(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
func (s *Server) OnClose(fn func()) {
	s.Server.RegisterOnShutdown(fn)
}
