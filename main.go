package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harlow/go-middleware-context/ctxhttp"
	"github.com/harlow/go-middleware-context/requestid"
	"github.com/harlow/go-middleware-context/userip"

	"github.com/harlow/go-middleware-context/Godeps/_workspace/src/golang.org/x/net/context"
)

type Server struct {
	context.Context
	ctxhttp.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(s.Context, w, r)
}

func requestIDMiddleware(next ctxhttp.Handler) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if reqID, ok := requestid.FromRequest(r); ok == nil {
			ctx = requestid.NewContext(ctx, reqID)
		}
		next.ServeHTTP(ctx, w, r)
	})
}

func userIPMiddleware(next ctxhttp.Handler) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if userIP, ok := userip.FromRequest(r); ok == nil {
			ctx = userip.NewContext(ctx, userIP)
		}
		next.ServeHTTP(ctx, w, r)
	})
}

func requestHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqID, _ := requestid.FromContext(ctx)
	userIP, _ := userip.FromContext(ctx)
	fmt.Fprintf(w, "Hello request: %s, from %s\n", reqID, userIP)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var handler ctxhttp.Handler
	handler = ctxhttp.HandlerFunc(requestHandler)
	handler = userIPMiddleware(handler)
	handler = requestIDMiddleware(handler)

	ctx := context.Background()
	svc := &Server{ctx, handler}
	log.Fatal(http.ListenAndServe(":"+port, svc))
}
