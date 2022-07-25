package mock

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

type MockHandler struct {
	W *httptest.ResponseRecorder
	R *http.Request
}

// Create mock objects
//
// Example:
//
// ```go
// mock := NewMock("", http.MethodGet, "/")
// ````
func NewMock(body string, method string, path string) *MockHandler {
	b := strings.NewReader(body)
	return NewMockByte(b, method, path)
}

// Create mock objects use bytes body
func NewMockByte(body io.Reader, method string, path string) *MockHandler {
	r := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()

	return &MockHandler{
		W: w,
		R: r,
	}
}

func NewGet(body string, path string) *MockHandler {
	return NewMock(body, http.MethodGet, path)
}

func NewPostJson(t *testing.T, body string, path string, data any) *MockHandler {
	mock := NewMock(body, http.MethodPost, path)
	mock.R.Header.Add("content-type", "application/json")

	b, err := json.Marshal(data)
	require.NoError(t, err)
	mock.R.Body = io.NopCloser(bytes.NewReader(b))

	return mock
}

func NewPostURLEncoded(body string, path string, data url.Values) *MockHandler {
	mock := NewMock(body, http.MethodPost, path)
	mock.R.Header.Add("content-type", "application/x-www-form-urlencoded")

	mock.R.PostForm = data

	return mock
}

func NewPostFormData(body string, path string, data contents.Multipart) *MockHandler {
	mock := NewMock(body, http.MethodPost, path)
	mock.R.Header.Add("content-type", data.ContentType())

	b := data.Export().Bytes()
	mock.R.Body = io.NopCloser(bytes.NewReader(b))

	return mock
}

// Set RemoteAddr
//
// default to 192.0.2.0/24 is "TEST-NET" in RFC 5737
func (c *MockHandler) SetAddr(addr string) {
	c.R.RemoteAddr = addr
}

// Set Host
//
// default to example.com
func (c *MockHandler) SetHost(host string) {
	c.R.Host = host
}

// Add handler
func (c *MockHandler) Handler(hand func(w http.ResponseWriter, r *http.Request)) {
	hand(c.W, c.R)
}

// Check if request success
func (c *MockHandler) Ok(t *testing.T) {
	c.Status(t, http.StatusOK)
}

// Check response status code
func (c *MockHandler) Status(t *testing.T, status int) {
	require.Equal(t, c.W.Code, status)
}

// Compare response body
func (c *MockHandler) EqBody(t *testing.T, body string) {
	require.Equal(t, c.W.Body.String(), body)
}

// Compare response body written json
func (c *MockHandler) EqJson(t *testing.T, obj any) {
	bytes, err := json.Marshal(obj)
	require.NoError(t, err)

	require.Equal(t, c.W.Body.Bytes(), bytes)
}
