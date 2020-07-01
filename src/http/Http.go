package http

import (
	"context"
	"net"
	"net/http"
)

//Server 新的http服务器,上下文使用tcp的服务器的
//Server 标准http
//ctx 上下文
type Server struct {
	Server *http.Server
	ctx    context.Context
}

func NewServer(ctx context.Context, addr string) *Server {
	ctx = context.WithValue(ctx, "test", "test")
	s := &Server{
		Server: &http.Server{Addr: addr, Handler: nil, MaxHeaderBytes: 1 << 30},
		ctx:    ctx,
	}
	s.Server.BaseContext = func(listener net.Listener) context.Context {
		return s.ctx
	}
	return s
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
