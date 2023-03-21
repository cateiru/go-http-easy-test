package easy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var client = new(http.Client)

type MockServer struct {
	Server *httptest.Server
	Header *http.Header

	Cookies []string
}

// Start mock server
func NewMockServer(handler http.Handler) *MockServer {
	server := httptest.NewServer(handler)

	return &MockServer{
		Server: server,
		Header: &http.Header{},
	}
}

// Start mock server with TLS mode
func NewMockTLSServer(handler http.Handler) *MockServer {
	server := httptest.NewTLSServer(handler)

	return &MockServer{
		Server: server,
	}
}

// close server
func (c *MockServer) Close() {
	c.Server.Close()
}

// convert to mock server url
func (c *MockServer) URL(path string) string {
	return c.Server.URL + path
}

func (c *MockServer) Cookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.Cookies = append(c.Cookies, cookie.String())
	}

	c.Header.Set("cookie", strings.Join(c.Cookies, "; "))
}

// GET Request
func (c *MockServer) Get(t *testing.T, path string) *Response {
	return c.Do(t, path, http.MethodGet, nil)
}

// Get Request and check status 200
func (c *MockServer) GetOK(t *testing.T, path string) *Response {
	resp := c.Get(t, path)
	resp.Ok(t)

	return resp
}

// POST Requests
func (c *MockServer) Post(t *testing.T, path string, contentType string, body io.Reader) *Response {
	c.Header.Add("Content-Type", contentType)

	return c.Do(t, path, http.MethodPost, body)
}

// application/x-www-form-urlencoded
func (c *MockServer) PostForm(t *testing.T, path string, value url.Values) *Response {
	return c.Post(t, path, "application/x-www-form-urlencoded", strings.NewReader(value.Encode()))
}

// application/json
func (c *MockServer) PostJson(t *testing.T, path string, obj any) *Response {
	b, err := json.Marshal(obj)
	require.NoError(t, err)

	return c.Post(t, path, "application/json", bytes.NewReader(b))
}

func (c *MockServer) PostString(t *testing.T, path string, contentType string, body string) *Response {
	r := strings.NewReader(body)
	resp := c.Post(t, path, contentType, r)

	return resp
}

// POST multipart/form-data
func (c *MockServer) PostFormData(t *testing.T, path string, form *Multipart) *Response {
	return c.FormData(t, path, http.MethodPost, form)
}

// multipart/form-data
func (c *MockServer) FormData(t *testing.T, path string, method string, form *Multipart) *Response {
	body := form.Export()

	c.Header.Add("Content-Type", form.ContentType())

	return c.Do(t, path, method, body)
}

func (c *MockServer) Do(t *testing.T, path string, method string, body io.Reader) *Response {
	r, err := http.NewRequest(method, c.URL(path), body)
	require.NoError(t, err)

	// insert headers
	for key, values := range *c.Header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	resp, err := client.Do(r)
	require.NoError(t, err)

	return NewResponse(resp)
}
