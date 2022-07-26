package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/contents"
	"github.com/stretchr/testify/require"
)

type MockHandler struct {
	W *httptest.ResponseRecorder
	R *http.Request
}

// Create mock objects.
// if set path is empty, replace to /.
//
// Example:
//
// 	m, err := NewMock("", http.MethodGet, "/")
func NewMock(body string, method string, path string) (*MockHandler, error) {
	b := strings.NewReader(body)
	return NewMockReader(b, method, path)
}

// Create mock objects use io.Reader body.
// if set path is empty, replace to /.
//
// Example:
//
// 	m, err := NewMockReader(strings.NewReader(""), http.MethodGet, "/")
func NewMockReader(body io.Reader, method string, path string) (*MockHandler, error) {
	// path case is `/`, `/name`, `https://` and `http://`
	reg := regexp.MustCompile(`(https?:\/)?\/.*`)

	if path == "" {
		path = "/"
	} else if !reg.MatchString(path) {
		return nil, errors.New("illegal path case. path: " + path)
	}

	r := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()

	return &MockHandler{
		W: w,
		R: r,
	}, nil
}

// GET Requests
func NewGet(body string, path string) (*MockHandler, error) {
	return NewMock(body, http.MethodGet, path)
}

// Post json. Use the POST or PUT method.
func NewJson(path string, data any, method string) (*MockHandler, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	mock, err := NewMockReader(bytes.NewReader(b), method, path)
	if err != nil {
		return nil, err
	}
	mock.R.Header.Add("content-type", "application/json")

	return mock, nil
}

// application/x-www-form-urlencoded type.
// Use the POST or PUT method.
func NewURLEncoded(path string, data url.Values, method string) (*MockHandler, error) {
	mock, err := NewMock(data.Encode(), method, path)
	if err != nil {
		return nil, err
	}
	mock.R.Header.Add("content-type", "application/x-www-form-urlencoded")

	return mock, nil
}

// multipart/form-data type.
// Use the POST or PUT method.
func NewFormData(path string, data *contents.Multipart, method string) (*MockHandler, error) {
	mock, err := NewMockReader(data.Export(), method, path)
	if err != nil {
		return nil, err
	}
	mock.R.Header.Add("content-type", data.ContentType())

	return mock, nil
}

// Set RemoteAddr
//
// default to 192.0.2.1:1234 is "TEST-NET" in RFC 5737
func (c *MockHandler) SetAddr(addr string) {
	c.R.RemoteAddr = addr
}

// Including cookies in the request
func (c *MockHandler) Cookie(cookies []*http.Cookie) {
	cookieLists := []string{}

	for _, cookie := range cookies {
		cookieLists = append(cookieLists, cookie.String())
	}

	c.R.Header.Set("cookie", strings.Join(cookieLists, "; "))
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
