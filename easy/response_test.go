package easy_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/v2/easy"
	"github.com/stretchr/testify/require"
)

func TestOk(t *testing.T) {
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}
	r := easy.NewResponse(resp)

	r.Ok(t)
}

func TestStatus(t *testing.T) {
	resp := &http.Response{
		StatusCode: 301,
	}
	r := easy.NewResponse(resp)

	r.Status(t, 301)
}

func TestEqBody(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,

		Body: io.NopCloser(strings.NewReader("aaaa")),
	}
	r := easy.NewResponse(resp)

	r.EqBody(t, "aaaa")
}

func TestEqJson(t *testing.T) {
	data := JsonData{
		Nya: "aaaa",
	}
	b, err := json.Marshal(data)
	require.NoError(t, err)

	resp := &http.Response{
		StatusCode: 200,

		Body: io.NopCloser(bytes.NewReader(b)),
	}
	r := easy.NewResponse(resp)

	r.EqJson(t, data)
}

func TestJson(t *testing.T) {
	data := JsonData{
		Nya: "aaaa",
	}
	b, err := json.Marshal(data)
	require.NoError(t, err)

	resp := &http.Response{
		StatusCode: 200,

		Body: io.NopCloser(bytes.NewReader(b)),
	}
	r := easy.NewResponse(resp)

	respBody := new(JsonData)
	err = r.Json(respBody)
	require.NoError(t, err)

	require.Equal(t, data, *respBody)
}

func TestSetCookies(t *testing.T) {
	c := &http.Cookie{
		Name:  "name",
		Value: "value",
	}

	resp := &http.Response{
		StatusCode: 200,

		Header: http.Header{
			"Set-Cookie": []string{c.String()},
		},
	}
	r := easy.NewResponse(resp)

	cookies := r.SetCookies()

	require.Len(t, cookies, 1)
	require.Equal(t, cookies[0].Name, c.Name)
	require.Equal(t, cookies[0].Value, c.Value)
}
