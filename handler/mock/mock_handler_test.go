package mock_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/contents"
	"github.com/cateiru/go-http-easy-test/handler/mock"
	"github.com/stretchr/testify/require"
)

type NewArgs struct {
	Body string
	Path string

	TestMessage string
}

type JsonData struct {
	Nya string `json:"nya"`
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

	t.Run("NewMockReader", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				m, err := mock.NewMockReader(strings.NewReader(c.Body), http.MethodGet, c.Path)
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
				_, err := mock.NewMockReader(strings.NewReader(c.Body), http.MethodGet, c.Path)
				require.Error(t, err)
			})
		}
	})

	t.Run("NewGet", func(t *testing.T) {
		m, err := mock.NewGet("", "/")
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		w := httptest.NewRecorder()

		require.Equal(t, m, &mock.MockHandler{
			R: r,
			W: w,
		})
	})

	t.Run("NewJson", func(t *testing.T) {
		data := JsonData{
			Nya: "hoge",
		}

		b, err := json.Marshal(data)
		require.NoError(t, err)

		t.Run("Success case", func(t *testing.T) {
			m, err := mock.NewJson("/", data, http.MethodPost)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			w := httptest.NewRecorder()

			r.Header.Set("content-type", "application/json")

			require.Equal(t, m, &mock.MockHandler{
				R: r,
				W: w,
			})

			require.Equal(t, m.R.Header.Get("content-type"), "application/json")
			body := new(JsonData)
			err = json.NewDecoder(m.R.Body).Decode(body)
			require.NoError(t, err)
			require.Equal(t, data, *body)
		})

		t.Run("empty body", func(t *testing.T) {
			m, err := mock.NewJson("/", "", http.MethodPost)
			require.NoError(t, err)

			buf := new(bytes.Buffer)
			buf.ReadFrom(m.R.Body)
			require.Equal(t, buf.String(), `""`)
		})
	})

	t.Run("NewPostURLEncoded", func(t *testing.T) {
		t.Run("Success get query", func(t *testing.T) {
			m, err := mock.NewURLEncoded("/", url.Values{"hoge": {"huga"}}, http.MethodPost)
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("hoge"), "huga")
		})

		t.Run("multi query", func(t *testing.T) {
			m, err := mock.NewURLEncoded("/", url.Values{"hoge": {"huga"}, "aaa": {"bbb"}}, http.MethodPost)
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("hoge"), "huga")
			require.Equal(t, m.R.FormValue("aaa"), "bbb")
		})

		t.Run("empty", func(t *testing.T) {
			m, err := mock.NewURLEncoded("/", url.Values{}, http.MethodPost)
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.Form, url.Values{})
		})
	})

	t.Run("NewPostFormData", func(t *testing.T) {
		t.Run("Success text from", func(t *testing.T) {
			multipart := contents.NewMultipart()
			err := multipart.Insert("key", "value")
			require.NoError(t, err)

			m, err := mock.NewFormData("/", multipart, http.MethodPost)
			require.NoError(t, err)

			require.Equal(t, m.R.Header.Get("content-type"), multipart.ContentType())

			err = m.R.ParseMultipartForm(32 << 20)
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("key"), "value")
		})

		t.Run("Success multi data", func(t *testing.T) {
			multipart := contents.NewMultipart()
			err := multipart.Insert("key", "value")
			require.NoError(t, err)
			err = multipart.Insert("aaa", "bbbb")
			require.NoError(t, err)

			m, err := mock.NewFormData("/", multipart, http.MethodPost)
			require.NoError(t, err)

			require.Equal(t, m.R.Header.Get("content-type"), multipart.ContentType())

			err = m.R.ParseMultipartForm(32 << 20)
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("key"), "value")
			require.Equal(t, m.R.FormValue("aaa"), "bbbb")
		})

		t.Run("Success file", func(t *testing.T) {
			file, err := os.Open("../../README.md")
			require.NoError(t, err)

			multipart := contents.NewMultipart()
			err = multipart.InsertFile("file", file)
			require.NoError(t, err)

			m, err := mock.NewFormData("/", multipart, http.MethodPost)
			require.NoError(t, err)

			require.Equal(t, m.R.Header.Get("content-type"), multipart.ContentType())

			err = m.R.ParseMultipartForm(32 << 20)
			require.NoError(t, err)

			formFile, _, err := m.R.FormFile("file")
			require.NoError(t, err)

			err = formFile.Close()
			require.NoError(t, err)
		})
	})
}

func TestSetAddr(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m, err := mock.NewMock("", http.MethodGet, "/")
		require.NoError(t, err)

		require.Equal(t, m.R.RemoteAddr, "192.0.2.1:1234")
	})

	t.Run("overwrite", func(t *testing.T) {
		m, err := mock.NewMock("", http.MethodGet, "/")
		require.NoError(t, err)

		m.SetAddr("203.0.113.0")

		require.Equal(t, m.R.RemoteAddr, "203.0.113.0")
	})
}

func TestCookie(t *testing.T) {
	cookie1 := http.Cookie{
		Name:  "session",
		Value: "12345",
	}

	cookie2 := http.Cookie{
		Name:  "aaaa",
		Value: "value",

		Secure:   true,
		HttpOnly: true,
	}

	t.Run("Success a cookie", func(t *testing.T) {
		m, err := mock.NewMock("", http.MethodGet, "/")
		require.NoError(t, err)

		m.Cookie([]*http.Cookie{
			&cookie1,
		})

		getCookie, err := m.R.Cookie("session")
		require.NoError(t, err)

		require.Equal(t, getCookie.Value, cookie1.Value)
	})

	t.Run("Success multi cookies", func(t *testing.T) {
		m, err := mock.NewMock("", http.MethodGet, "/")
		require.NoError(t, err)

		m.Cookie([]*http.Cookie{
			&cookie1,
			&cookie2,
		})

		t.Log(cookie1.String())

		getCookie1, err := m.R.Cookie("session")
		require.NoError(t, err)

		require.Equal(t, getCookie1.Value, cookie1.Value)

		getCookie2, err := m.R.Cookie("aaaa")
		require.NoError(t, err)

		require.Equal(t, getCookie2.Value, cookie2.Value)
	})
}
