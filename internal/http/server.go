package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(handler http.Handler, port string) *Server {
    return &Server{
        httpServer: &http.Server{
            Addr:         ":" + port,
            Handler:      handler,
            ReadTimeout:  5 * time.Second,
            WriteTimeout: 10 * time.Second,
            IdleTimeout:  60 * time.Second,
        },
    }
}

func (s *Server) Start() error {
    log.Println("Server started on", s.httpServer.Addr)
    return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}