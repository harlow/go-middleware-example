package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/harlow/go-middleware-context/ctxhttp"
	"github.com/harlow/go-middleware-context/requestid"
	"github.com/harlow/go-middleware-context/userip"

	"golang.org/x/net/context"
)

func requestIDMiddleware(next ctxhttp.Handler) ctxhttp.Handler {
	return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		requestID := requestid.FromRequest(r)
		ctx = requestid.NewContext(ctx, requestID)
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
	var handler ctxhttp.Handler
	handler = ctxhttp.HandlerFunc(requestHandler)
	handler = userIPMiddleware(handler)
	handler = requestIDMiddleware(handler)

	ctx := context.Background()
	svc := &ctxhttp.Server{ctx, handler}
	log.Fatal(http.ListenAndServe(":8080", svc))
}
