# HTTP Server w/ Context

Example of using `context` package with HTTP Middleware in Golang

## Usage

Run the server:

    $ go run main.go

In another terminal cURL the endpoint:

    $ curl localhost:8080 -H "X-Request-Id: UNIQUE_REQ_ID"
    Hello request: UNIQUE_REQ_ID, from 127.0.0.1

## Credits

This codebase was heavily inspired by the following talks and repositories:

* https://github.com/alexedwards/stack
* https://joeshaw.org/net-context-and-http-handler/
