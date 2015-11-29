# HTTP Server w/ Context

Example of using [package context][1] with HTTP Middleware in Golang

## Usage

Run the server:

    $ go run main.go

In another terminal cURL the endpoint:

    $ curl localhost:8080 -H "X-Request-Id: UNIQUE_REQ_ID"
    Hello request: UNIQUE_REQ_ID, from 127.0.0.1

## Credits

This codebase was heavily inspired by the following blog posts and repositories:

* http://www.alexedwards.net/blog/making-and-using-middleware
* https://joeshaw.org/net-context-and-http-handler/
* https://github.com/alexedwards/stack

[1]: https://godoc.org/golang.org/x/net/context
