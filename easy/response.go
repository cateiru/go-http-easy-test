package easy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Response struct {
	Resp *http.Response
}

func NewResponse(resp *http.Response) *Response {
	return &Response{
		Resp: resp,
	}
}

func (c *Response) Ok(t *testing.T) {
	c.Status(t, 200)
}

func (c *Response) Status(t *testing.T, status int) {
	require.Equal(t, c.Resp.StatusCode, status)
}

func (c *Response) Body() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Resp.Body)

	return buf
}

func (c *Response) EqBody(t *testing.T, body string) {
	require.Equal(t, body, c.Body().String())
}

func (c *Response) EqJson(t *testing.T, obj any) {
	bytes, err := json.Marshal(obj)
	require.NoError(t, err)

	require.Equal(t, c.Body().Bytes(), bytes)
}

// Prase json body
func (c *Response) Json(v any) error {
	data := c.Body().Bytes()
	return json.Unmarshal(data, v)
}

// Returns set-cookie's
func (c *Response) SetCookies() []*http.Cookie {
	return c.Resp.Cookies()
}
