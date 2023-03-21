package easy

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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type MockHandler struct {
	W *httptest.ResponseRecorder
	R *http.Request

	Cookies []string
}

// Create mock objects.
// if set path is empty, replace to /.
//
// Example:
//
//	m, err := NewMock("", http.MethodGet, "/")
func NewMock(path string, method string, body string) (*MockHandler, error) {
	b := strings.NewReader(body)
	return NewMockReader(path, method, b)
}

// Create mock objects use io.Reader body.
// if set path is empty, replace to /.
//
// Example:
//
//	m, err := NewMockReader(strings.NewReader(""), http.MethodGet, "/")
func NewMockReader(path string, method string, body io.Reader) (*MockHandler, error) {
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

// Post json. Use the POST or PUT method.
func NewJson(path string, method string, data any) (*MockHandler, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	mock, err := NewMockReader(path, method, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	mock.R.Header.Add("content-type", "application/json")

	return mock, nil
}

// application/x-www-form-urlencoded type.
// Use the POST or PUT method.
func NewURLEncoded(path string, method string, data url.Values) (*MockHandler, error) {
	mock, err := NewMock(path, method, data.Encode())
	if err != nil {
		return nil, err
	}
	mock.R.Header.Add("content-type", "application/x-www-form-urlencoded")

	return mock, nil
}

// multipart/form-data type.
// Use the POST or PUT method.
func NewFormData(path string, method string, data *Multipart) (*MockHandler, error) {
	mock, err := NewMockReader(path, method, data.Export())
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

	for _, cookie := range cookies {
		c.Cookies = append(c.Cookies, cookie.String())
	}

	c.R.Header.Set("cookie", strings.Join(c.Cookies, "; "))
}

// Add handler
func (c *MockHandler) Handler(hand func(w http.ResponseWriter, r *http.Request)) {
	hand(c.W, c.R)
}

// Returns echo context
// this method is require labstack/echo package.
//
// Usage:
//
//	c := m.Echo()
//	c.SetPath("/users/:email")
//	c.SetParamNames("email")
//	c.SetParamValues("jon@labstack.com")
func (c *MockHandler) Echo() echo.Context {
	e := echo.New()
	return e.NewContext(c.R, c.W)
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

// Prase json body
func (c *MockHandler) Json(v any) error {
	data := c.W.Body.Bytes()
	return json.Unmarshal(data, v)
}

// Returns Set-Cookie headers
func (c *MockHandler) SetCookies() []*http.Cookie {
	return c.Response().Cookies()
}

// Returns response
func (c *MockHandler) Response() *http.Response {
	return c.W.Result()
}

// Returns set-cookie
func (c *MockHandler) FindCookie(name string) *http.Cookie {
	cookies := c.Response().Cookies()

	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
