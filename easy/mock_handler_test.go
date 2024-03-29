package easy_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/v2/easy"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type NewArgs struct {
	Body string
	Path string

	TestMessage string
}

func EchoHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
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
				m, err := easy.NewMock(c.Path, http.MethodGet, c.Body)
				require.NoError(t, err)

				// overwrite
				if c.Path == "" {
					c.Path = "/"
				}

				r := httptest.NewRequest(http.MethodGet, c.Path, strings.NewReader(c.Body))
				w := httptest.NewRecorder()

				require.Equal(t, m, &easy.MockHandler{
					R: r,
					W: w,
				}, c.TestMessage)
			})
		}

		for _, c := range failedCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				_, err := easy.NewMock(c.Path, http.MethodGet, c.Body)
				require.Error(t, err)
			})
		}
	})

	t.Run("NewMockReader", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				m, err := easy.NewMockReader(c.Path, http.MethodGet, strings.NewReader(c.Body))
				require.NoError(t, err)

				// overwrite
				if c.Path == "" {
					c.Path = "/"
				}

				r := httptest.NewRequest(http.MethodGet, c.Path, strings.NewReader(c.Body))
				w := httptest.NewRecorder()

				require.Equal(t, m, &easy.MockHandler{
					R: r,
					W: w,
				}, c.TestMessage)
			})
		}

		for _, c := range failedCases {
			t.Run(c.TestMessage, func(t *testing.T) {
				_, err := easy.NewMockReader(c.Path, http.MethodGet, strings.NewReader(c.Body))
				require.Error(t, err)
			})
		}
	})

	t.Run("NewJson", func(t *testing.T) {
		data := JsonData{
			Nya: "hoge",
		}

		b, err := json.Marshal(data)
		require.NoError(t, err)

		t.Run("Success case", func(t *testing.T) {
			m, err := easy.NewJson("/", http.MethodPost, data)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			w := httptest.NewRecorder()

			r.Header.Set("content-type", "application/json")

			require.Equal(t, m, &easy.MockHandler{
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
			m, err := easy.NewJson("/", http.MethodPost, "")
			require.NoError(t, err)

			buf := new(bytes.Buffer)
			buf.ReadFrom(m.R.Body)
			require.Equal(t, buf.String(), `""`)
		})
	})

	t.Run("NewPostURLEncoded", func(t *testing.T) {
		t.Run("Success get query", func(t *testing.T) {
			m, err := easy.NewURLEncoded("/", http.MethodPost, url.Values{"hoge": {"huga"}})
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("hoge"), "huga")
		})

		t.Run("multi query", func(t *testing.T) {
			m, err := easy.NewURLEncoded("/", http.MethodPost, url.Values{"hoge": {"huga"}, "aaa": {"bbb"}})
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("hoge"), "huga")
			require.Equal(t, m.R.FormValue("aaa"), "bbb")
		})

		t.Run("empty", func(t *testing.T) {
			m, err := easy.NewURLEncoded("/", http.MethodPost, url.Values{})
			require.NoError(t, err)

			err = m.R.ParseForm()
			require.NoError(t, err)

			require.Equal(t, m.R.Form, url.Values{})
		})
	})

	t.Run("NewPostFormData", func(t *testing.T) {
		t.Run("Success text from", func(t *testing.T) {
			multipart := easy.NewMultipart()
			err := multipart.Insert("key", "value")
			require.NoError(t, err)

			m, err := easy.NewFormData("/", http.MethodPost, multipart)
			require.NoError(t, err)

			require.Equal(t, m.R.Header.Get("content-type"), multipart.ContentType())

			err = m.R.ParseMultipartForm(32 << 20)
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("key"), "value")
		})

		t.Run("Success multi data", func(t *testing.T) {
			multipart := easy.NewMultipart()
			err := multipart.Insert("key", "value")
			require.NoError(t, err)
			err = multipart.Insert("aaa", "bbbb")
			require.NoError(t, err)

			m, err := easy.NewFormData("/", http.MethodPost, multipart)
			require.NoError(t, err)

			require.Equal(t, m.R.Header.Get("content-type"), multipart.ContentType())

			err = m.R.ParseMultipartForm(32 << 20)
			require.NoError(t, err)

			require.Equal(t, m.R.FormValue("key"), "value")
			require.Equal(t, m.R.FormValue("aaa"), "bbbb")
		})

		t.Run("Success file", func(t *testing.T) {
			file, err := os.Open("../README.md")
			require.NoError(t, err)

			multipart := easy.NewMultipart()
			err = multipart.InsertFile("file", file)
			require.NoError(t, err)

			m, err := easy.NewFormData("/", http.MethodPost, multipart)
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
		m, err := easy.NewMock("/", http.MethodGet, "")
		require.NoError(t, err)

		require.Equal(t, m.R.RemoteAddr, "192.0.2.1:1234")
	})

	t.Run("overwrite", func(t *testing.T) {
		m, err := easy.NewMock("/", http.MethodGet, "")
		require.NoError(t, err)

		m.SetAddr("203.0.113.0")

		require.Equal(t, m.R.RemoteAddr, "203.0.113.0")
	})
}

func TestMockCookie(t *testing.T) {
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
		m, err := easy.NewMock("/", http.MethodGet, "")
		require.NoError(t, err)

		m.Cookie([]*http.Cookie{
			&cookie1,
		})

		getCookie, err := m.R.Cookie("session")
		require.NoError(t, err)

		require.Equal(t, getCookie.Value, cookie1.Value)
	})

	t.Run("Success multi cookies", func(t *testing.T) {
		m, err := easy.NewMock("/", http.MethodGet, "")
		require.NoError(t, err)

		m.Cookie([]*http.Cookie{
			&cookie1,
			&cookie2,
		})

		getCookie1, err := m.R.Cookie("session")
		require.NoError(t, err)

		require.Equal(t, getCookie1.Value, cookie1.Value)

		getCookie2, err := m.R.Cookie("aaaa")
		require.NoError(t, err)

		require.Equal(t, getCookie2.Value, cookie2.Value)
	})
}

func TestHandler(t *testing.T) {
	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(Handler)

	require.Equal(t, m.W.Body.String(), "OK")
}

func TestEcho(t *testing.T) {
	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	c := m.Echo()

	require.NoError(t, EchoHandler(c))
}

func TestMockOk(t *testing.T) {
	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(Handler)

	m.Ok(t)
}

func TestMockEqBody(t *testing.T) {
	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(Handler)

	m.EqBody(t, "OK")
}

func TestMockEqJson(t *testing.T) {
	data := JsonData{
		Nya: "aaaa",
	}
	b, err := json.Marshal(data)
	require.NoError(t, err)

	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})

	m.EqJson(t, data)
}

func TestMockJson(t *testing.T) {
	data := JsonData{
		Nya: "aaaa",
	}
	b, err := json.Marshal(data)
	require.NoError(t, err)

	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})

	resp := new(JsonData)
	err = m.Json(resp)
	require.NoError(t, err)

	require.Equal(t, data, *resp)
}

func TestMockSetCookies(t *testing.T) {
	c := &http.Cookie{
		Name:  "name",
		Value: "value",
	}

	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, c)
	})

	cookies := m.SetCookies()

	require.Len(t, cookies, 1)
	require.Equal(t, cookies[0].Name, c.Name)
	require.Equal(t, cookies[0].Value, c.Value)
}

func TestResponse(t *testing.T) {
	m, err := easy.NewMock("/", http.MethodGet, "")
	require.NoError(t, err)

	m.Handler(func(w http.ResponseWriter, r *http.Request) {
	})

	resp := m.Response()

	require.NotNil(t, resp)
}
