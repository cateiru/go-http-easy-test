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
go get -u github.com/cateiru/go-http-easy-test/v2
```

## Mock

The user can choose from the following two options.

- Actually start the server using `httptest.NewServer`
- Mock the Handler arguments (`w http.ResponseWriter, r *http.Request`)

### Actually start the server using `httptest.NewServer`

```go
package main_test

import (
    "testing"

    "github.com/cateiru/go-http-easy-test/easy"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    ...do something
}

func TestHandler(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", Handler)

    // create server
    s := easy.NewMockServer(mux)
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

    // Easily build multipart/form-data
    form := easy.NewMultipart()
    form.Insert("key", "value")
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

    // prase response json
    body := new(JsonType)
    err := resp.Json(body)

    // returns Set-Cookie headers
    cookies := resp.SetCookies()
}
```

### Mock the Handler arguments (`w http.ResponseWriter, r *http.Request`)

```go
package main_test

import (
    "testing"

    "github.com/cateiru/go-http-easy-test/easy"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    ...do something
}

func EchoHandler(c echo.Context) error {
    ...do something
}

func TestHandler(t *testing.T) {
    // Default
    m, err := easy.NewMock(body, http.MethodGet, "/")
    m, err := easy.NewMockReader(reader, http.MethodGet, "/")

    // GET
    m, err := easy.NewGet(body, "/")

    // POST or PUT send json
    m, err := easy.NewJson("/", data, http.MethodPost)

    // POST or PUT send x-www-form-urlencoded
    m, err := easy.NewURLEncoded("/", url, http.MethodPost)

    // POST or PUT send multipart/form-data
    // Easily build multipart/form-data using the `contents` package.
    m, err := easy.NewFormData("/", multipart, http.MethodPost)


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

        // prase response json
    body := new(JsonType)
    err := m.Json(body)

    // returns Set-Cookie headers
    cookies := m.SetCookies()

    // Return http.Response
    response := m.Response()

    // set-cookie
    cookie := m.FindCookie("name")
}
```

### multipart

Easily create `multipart/form-data` requests.<br/>
This method is used when submitting with `multipart/form-data`.

```go
package main

import (
    "os"

    "github.com/cateiru/go-http-easy-test/easy"
)


func main() {
    m := easy.NewMultipart()

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
