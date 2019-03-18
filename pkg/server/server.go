package server

import (
	"context"
	"net/http"

	"github.com/dchertkov/scrapper/pkg/config"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func (s *Server) Listen(ctx context.Context) (err error) {
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	err = srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return
}

func NewServer(conf config.Server, h http.Handler) *Server {
	return &Server{
		Addr:    conf.Addr(),
		Handler: h,
	}
}
