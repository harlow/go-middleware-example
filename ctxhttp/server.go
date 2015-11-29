package ctxhttp

import (
	"net/http"

	"golang.org/x/net/context"
)

type Server struct {
	context.Context
	Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(s.Context, w, r)
}
