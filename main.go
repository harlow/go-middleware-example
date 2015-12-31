package main

import (
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/harlow/go-middleware-context/ctxhttp"
	"github.com/harlow/go-middleware-context/requestid"
	"github.com/harlow/go-middleware-context/userip"

	"github.com/paulbellamy/ratecounter"
	"golang.org/x/net/context"
)

var (
  counter        *ratecounter.RateCounter
  hitsperminute = expvar.NewInt("hits_per_minute")
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

func requestCtrMiddleware(next ctxhttp.Handler) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		counter.Incr(1)
		hitsperminute.Set(counter.Rate())
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
	counter = ratecounter.NewRateCounter(1 * time.Minute)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var handler ctxhttp.Handler
	handler = ctxhttp.HandlerFunc(requestHandler)
	handler = userIPMiddleware(handler)
	handler = requestIDMiddleware(handler)
	handler = requestCtrMiddleware(handler)

	ctx := context.Background()
	svc := &Server{ctx, handler}

	http.HandleFunc("/", svc.ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
