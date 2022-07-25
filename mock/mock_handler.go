package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockHandler struct {
	W *httptest.ResponseRecorder
	R *http.Request
}

// handlerのwとrをmockする
func NewMock(body string, method string, path string) *MockHandler {
	b := strings.NewReader(body)
	return NewMockByte(b, method, path)
}

func NewMockByte(body io.Reader, method string, path string) *MockHandler {
	r := httptest.NewRequest(method, path, body)

	w := httptest.NewRecorder()

	return &MockHandler{
		W: w,
		R: r,
	}
}

func (c *MockHandler) MockHandler(hand func(w http.ResponseWriter, r *http.Request)) {
	hand(c.W, c.R)
}

func (c *MockHandler) Ok(t *testing.T) {
	c.Status(t, http.StatusOK)
}

func (c *MockHandler) Status(t *testing.T, status int) {
	require.Equal(t, c.W.Code, status)
}

func (c *MockHandler) EqBody(t *testing.T, body string) {
	require.Equal(t, c.W.Body.String(), body)
}
