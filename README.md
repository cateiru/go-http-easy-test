# Go http easy test

Goのhttpテストを楽するライブラリ

## Usage

```bash
go get -u github.com/cateiru/go-http-easy-test
```

### Handlerのwとrをモックしてテストする

```go
func SampleHandler(w http.ResponseWriter, r *http.Request) {
    ...
}

func TestHandler(t *testing.T) {
    requestBody := ""
    requestMethod := http.MethodGet
    requestPath := "/"

    h := mock.NewMock(requestBody, requestMethod, requestPath)

    h.MockHandler(SampleHandler)

    // Test if the response code is 200
    h.Ok(t)

    // custom status code
    h.Status(t, 403)

    // Test body
    testResponseBody := "response"
    h.EqBody(t, testResponseBOdy)
}

```

### サーバを立ててE2Eテストをする

```go
func SampleHandler(w http.ResponseWriter, r *http.Request) {
    ...
}

func TestServer(t *testing.T) {
    s := mock.NewMockServer(SampleHandler)
    defer s.Close()

    // Get
    resp := s.Get(t, "/")

    // Get and check response status code
    resp := s.GetOk(t, "/")

    // Post
    body := strings.NewReader(`{"key": "value"}`)
    resp := s.Post(t, "/", "application/json", body)

    // Post application/x-www-form-urlencoded
    url := url.Values{
        "aaa": {"aaaa"}
    }
    resp := s.PostForm(t, "/", url)

    // Post str
    body := `{"key": "value"}`
    resp := s.PostString(t, "/", "application/json", body)

    // Post multipart/form-data
    form := contents.NewMultipart()
    form.Insert(t, "key", "value")
    form.InsertFile(t, "key", file)

    resp := s.PostFormData(t, "/", form)

    // other method version
    resp := s.FormData(t, "/", http.MethodPut, form)

    // other
    resp := s.Do(t, "/", http.MethodDelete, body, func (r *http.Request) {
        r.Header.Add("Authentication", "Basic XXXXXX:XXXXXX")
    })
}
```

## License

[MIT](./LICENSE)
