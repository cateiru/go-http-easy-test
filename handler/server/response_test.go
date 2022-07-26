package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/handler/server"
	"github.com/stretchr/testify/require"
)

func TestOk(t *testing.T) {
	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}
	r := server.NewResponse(resp)

	r.Ok(t)
}

func TestStatus(t *testing.T) {
	resp := &http.Response{
		StatusCode: 301,
	}
	r := server.NewResponse(resp)

	r.Status(t, 301)
}

func TestEqBody(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,

		Body: io.NopCloser(strings.NewReader("aaaa")),
	}
	r := server.NewResponse(resp)

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
	r := server.NewResponse(resp)

	r.EqJson(t, data)
}
