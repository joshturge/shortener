package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joshturge/shortener/pkg/repo"
	"golang.org/x/sync/errgroup"
)

// Server will register routes and serve this web application
type Server struct {
	srv *http.Server
	rep repo.Repository
}

// NewServer will create a new server given the address to listen on and a repository to store urls
func NewServer(addr string, r repo.Repository) *Server {

	// register routes
	router := NewRouter()
	router.HandlerFunc(`^/$|/favicon\.ico`, Index)
	router.HandlerFunc(`^/[0-9A-Fa-f]{8}$`, Url(r))
	router.HandlerFunc("/shorten", Shorten(r))

	return &Server{&http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, r}
}

// Serve will listen and serve this app then exit when a signal has been received
func (s *Server) Serve(sig chan os.Signal) (err error) {
	go func() {
		if err = s.srv.ListenAndServe(); err != nil {
			err = fmt.Errorf("encountered an error while listen and serve: %w", err)
		}
	}()

	<-sig

	return err
}

// Close will call the servers close method and return any errors
func (s *Server) Close() error {
	errg, _ := errgroup.WithContext(context.Background())
	errg.Go(s.rep.Close)
	errg.Go(s.srv.Close)

	return errg.Wait()
}
