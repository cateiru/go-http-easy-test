package mock_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/handler/mock"
	"github.com/stretchr/testify/require"
)

type NewArgs struct {
	Body string
	Path string

	TestMessage string
}

func TestNewMock(t *testing.T) {
	successCases := []NewArgs{
		{
			Body:        "hogehoge",
			Path:        "/",
			TestMessage: "default case",
		},
		{
			Body:        "hogehoge",
			Path:        "/aaaaa",
			TestMessage: "no root path case",
		},
		{
			Body:        "hogehoge",
			Path:        "/aaaaa?hoge=huga",
			TestMessage: "no root path and exists query param case",
		},
		{
			Body:        "",
			Path:        "/",
			TestMessage: "empty body case",
		},
		{
			Body:        "aaa",
			Path:        "",
			TestMessage: "Empty path case",
		},
		{
			Body:        "aaa",
			Path:        "https://cateiru.com/",
			TestMessage: "exists host in path case",
		},
		{
			Body:        "aaa",
			Path:        "http://cateiru.com/",
			TestMessage: "exists host and no TLS in path case",
		},
		{
			Body:        "aaa",
			Path:        "https://cateiru.com/aaaaa",
			TestMessage: "exists host and no root path case",
		},
		{
			Body:        "aaa",
			Path:        "https://cateiru.com/aaaaa?hoge=huga",
			TestMessage: "exists host, no root path and query param case",
		},
	}

	failedCases := []NewArgs{
		{
			Body:        "",
			Path:        "aaaaa",
			TestMessage: "illegal path case",
		},
	}

	t.Run("NewMock", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				m, err := mock.NewMock(c.Body, http.MethodGet, c.Path)
				require.NoError(t, err)

				// overwrite
				if c.Path == "" {
					c.Path = "/"
				}

				r := httptest.NewRequest(http.MethodGet, c.Path, strings.NewReader(c.Body))
				w := httptest.NewRecorder()

				require.Equal(t, m, &mock.MockHandler{
					R: r,
					W: w,
				}, c.TestMessage)
			})
		}

		for _, c := range failedCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				_, err := mock.NewMock(c.Body, http.MethodGet, c.Path)
				require.Error(t, err)
			})
		}
	})
}
