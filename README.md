# Go http easy test

[![Go Reference](https://pkg.go.dev/badge/github.com/cateiru/go-http-easy-test.svg)](https://pkg.go.dev/github.com/cateiru/go-http-easy-test) [![Go Report Card](https://goreportcard.com/badge/github.com/cateiru/go-http-easy-test)](https://goreportcard.com/report/github.com/cateiru/go-http-easy-test) [![Go](https://github.com/cateiru/go-http-easy-test/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/cateiru/go-http-easy-test/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/cateiru/go-http-easy-test/branch/main/graph/badge.svg?token=3yN9nRKyvb)](https://codecov.io/gh/cateiru/go-http-easy-test)

A package that wraps `net/http/httptest` and allows you to easily test HTTP Handlers.

✅ Easy<br/>
✅ Intuitive<br/>
✅ Support `application/json`<br/>
✅ Support `application/x-www-form-urlencoded`<br/>
✅ Support `multipart/form-data`<br/>
✅ Support [Echo package](https://echo.labstack.com/)<br/>
✅ Support cookie<br/>

## Install

```bash
go get -u github.com/cateiru/go-http-easy-test
```

## `handler` package

The `handler` package provides methods that make it easier to test HTTP Handlers.<br/>
The user can choose from the following two options.

- Actually start the server using `httptest.NewServer`
- Mock the Handler arguments (`w http.ResponseWriter, r *http.Request`)

### Actually start the server using `httptest.NewServer`

```go
package main_test

import (
    "testing"

    "github.com/cateiru/go-http-easy-test/handler/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    ...do something
}

func TestHandler(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", Handler)

    // create server
    s := server.NewMockServer(mux)
    // Start the server with TLS using:
    // s := server.TestNewMockTLSServer(mux)
    defer s.Close()

    // Option: You can set cookies.
    cookie := &http.Cookie{
        Name:  "name",
        Value: "value",
    }
    s.Cookie([]*http.Cookie{
        cookie,
    })

    // GET
    resp := s.Get(t, "/")
    resp := s.GetOK(t, "/")

    // POST
    resp := s.Post(t, "/", "text/plain", body)
    resp := s.PostForm(t, "/", url) // application/x-www-form-urlencoded
    resp := s.PostJson(t, "/", obj) // application/json
    resp := s.PostString(t, "/", "text/plain", body)

    // Easily build multipart/form-data using the `contents` package.
    resp := s.PostFormData(t, "/", form)
    resp := s.FormData(t, "/", "[method]", form)

    // Other
    resp := s.Do(t, "/", "[method]", body)

    // The `resp` of all return values are easy to compare.
    // Check status
    resp.Ok(t)
    resp.Status(t, 200)

    // get body
    body := resp.Body().String()

    // Compare response body
    resp.EqBody(t, body)
    resp.EqJson(t, obj)
}
```

### Mock the Handler arguments (`w http.ResponseWriter, r *http.Request`)

```go
package main_test

import (
    "testing"

    "github.com/cateiru/go-http-easy-test/handler/mock"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    ...do something
}

func EchoHandler(c echo.Context) error {
    ...do something
}

func TestHandler(t *testing.T) {
    // Default
    m, err := mock.NewMock(body, http.MethodGet, "/")
    m, err := mock.NewMockReader(reader, http.MethodGet, "/")

    // GET
    m, err := mock.NewGet(body, "/")

    // POST or PUT send json
    m, err := mock.NewJson("/", data, http.MethodPost)

    // POST or PUT send x-www-form-urlencoded
    m, err := mock.NewURLEncoded("/", url, http.MethodPost)

    // POST or PUT send multipart/form-data
    // Easily build multipart/form-data using the `contents` package.
    m, err := mock.NewFormData("/", multipart, http.MethodPost)


    // Option: set remote addr
    m.SetAddr("203.0.113.0")

    // Option: You can set cookies.
    cookie := &http.Cookie{
        Name:  "name",
        Value: "value",
    }
    m.Cookie([]*http.Cookie{
        cookie,
    })

    // Set handler and run
    m.Handler(Handler)

    // Use echo package
    echoCtx := m.Echo()
    err := EchoHandler(echoCtx)

    // check response
    m.Ok(t)
    m.Status(t, 200)

    // Compare response body
    m.EqBody(t, body)
    m.EqJson(t, obj)
}
```

## `contents` package

The `contents` package provides methods that simplify HTTP-related testing.

### multipart

Easily create `multipart/form-data` requests.<br/>
This method is also used when submitting with `multipart/form-data` using the `handler` package.

```go
package main

import (
    "os"

    "github.com/cateiru/go-http-easy-test/contents"
)


func main() {
    m := contents.NewMultipart()

    // Add a string format form.
    err := m.Insert("key", "value")

    // Add a file format form.
    file, err := os.Open("path")
    err := m.InsertFile("key", file)

    // Outputs in the specified format.
    body := m.Export()
    contentType := m.ContentType()

    // Use `handler` package
    // Actually start the server using `httptest.NewServer`
    s := server.NewMockServer(mux)
    defer s.Close()
    resp := s.PostFormData(t, "/", m)
    // Mock the Handler arguments (`w http.ResponseWriter, r *http.Request`)
    m, err := mock.NewFormData("/", m, http.MethodPost)
}

```

## License

[MIT](./LICENSE)
