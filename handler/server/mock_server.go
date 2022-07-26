package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/contents"
	"github.com/stretchr/testify/require"
)

var client = new(http.Client)

type MockServer struct {
	Server *httptest.Server
}

// モックサーバーを起動する
func NewMockServer(handler http.Handler) *MockServer {
	server := httptest.NewServer(handler)

	return &MockServer{
		Server: server,
	}
}

// TLSでモックサーバーを起動する
func NewMockTLSServer(handler http.Handler) *MockServer {
	server := httptest.NewTLSServer(handler)

	return &MockServer{
		Server: server,
	}
}

// サーバーをclose
func (c *MockServer) Close() {
	c.Server.Close()
}

// モックサーバー用のURLに変換する
func (c *MockServer) URL(path string) string {
	return c.Server.URL + path
}

// Getメソッドで取得
func (c *MockServer) Get(t *testing.T, path string) *Response {
	resp, err := http.Get(c.URL(path))
	require.NoError(t, err)

	return NewResponse(resp)
}

// Getメソッドで取得し、レスポンスステータスが200かどうかを確認する
func (c *MockServer) GetOK(t *testing.T, path string) *Response {
	resp := c.Get(t, path)
	resp.Ok(t)

	return resp
}

func (c *MockServer) Post(t *testing.T, path string, contentType string, body io.Reader) *Response {
	resp, err := http.Post(c.URL(path), contentType, body)
	require.NoError(t, err)

	return NewResponse(resp)
}

// application/x-www-form-urlencoded
func (c *MockServer) PostForm(t *testing.T, path string, value url.Values) *Response {
	resp, err := http.PostForm(c.URL(path), value)
	require.NoError(t, err)

	return NewResponse(resp)
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
func (c *MockServer) PostFormData(t *testing.T, path string, form *contents.Multipart) *Response {
	return c.FormData(t, path, http.MethodPost, form)
}

// multipart/form-data
func (c *MockServer) FormData(t *testing.T, path string, method string, form *contents.Multipart) *Response {
	body := form.Export()

	return c.Do(t, path, method, body, func(r *http.Request) {
		r.Header.Add("Content-Type", form.ContentType())
	})
}

func (c *MockServer) Do(t *testing.T, path string, method string, body io.Reader, before func(r *http.Request)) *Response {
	r, err := http.NewRequest(method, c.URL(path), body)
	require.NoError(t, err)

	before(r)

	resp, err := client.Do(r)
	require.NoError(t, err)

	return NewResponse(resp)
}
