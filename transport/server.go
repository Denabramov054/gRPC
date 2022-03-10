package transport

import (
	"google.golang.org/grpc"
	"grpc/pkg/api"
	"grpc/transport/handler"
	"net"
)

type Server struct {
	srv *grpc.Server
	Deps Deps
}

type Deps struct {
 UserHandler *handler.UserHandler

}

type ServerConfig struct {
	Host string
	Port string
}

func NewServer(d Deps) *Server {
	return &Server{
		srv: grpc.NewServer(),
		Deps: d,
	}
}

func(s *Server) ListenAndServe(cfg ServerConfig) error {
	api.RegisterUserServer(s.srv, s.Deps.UserHandler)
	l, err := net.Listen("tcp", cfg.Host+":"+cfg.Port)
	if err != nil {
		return err
	}
	if err = s.srv.Serve(l); err !=nil {
		return err
	}
	return nil
}