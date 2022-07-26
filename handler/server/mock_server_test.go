package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cateiru/go-http-easy-test/contents"
	"github.com/cateiru/go-http-easy-test/handler/server"
	"github.com/stretchr/testify/require"
)

type JsonData struct {
	Nya string `json:"nya"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func TestNewMockServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	s := server.NewMockServer(mux)
	defer s.Close()
	require.Regexp(t, `http:\/\/.+`, s.Server.URL)
}

func TestNewMockTLSServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	s := server.NewMockTLSServer(mux)
	defer s.Close()
	require.Regexp(t, `https:\/\/.+`, s.Server.URL)
}

func TestURL(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	s := server.NewMockServer(mux)
	defer s.Close()

	url := s.URL("/aaaaa")

	require.Regexp(t, `http:\/\/.+/aaaaa`, url)
}

func TestGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	s := server.NewMockServer(mux)
	defer s.Close()

	resp := s.Get(t, "/")

	resp.Ok(t)

	b := resp.Body().String()
	require.Equal(t, b, "OK")
}

func TestGetOk(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	s := server.NewMockServer(mux)
	defer s.Close()

	s.GetOK(t, "/")

}

func TestPost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		w.Write(buf.Bytes())
	})

	s := server.NewMockServer(mux)
	defer s.Close()

	body := "hello"

	resp := s.Post(t, "/", "plain/text", strings.NewReader(body))
	resp.Ok(t)

	require.Equal(t, body, resp.Body().String())
}

func TestPostFrom(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		value := r.PostFormValue("key")

		w.Write([]byte(value))
	})

	s := server.NewMockServer(mux)
	defer s.Close()

	data := url.Values{"key": {"aaaa"}}

	resp := s.PostForm(t, "/", data)
	resp.Ok(t)

	require.Equal(t, data.Get("key"), resp.Body().String())
}

func TestPostJson(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		obj := new(JsonData)
		err := json.Unmarshal(buf.Bytes(), obj)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		w.Write([]byte(obj.Nya))
	})

	s := server.NewMockServer(mux)
	defer s.Close()

	data := JsonData{
		Nya: "aaaaa",
	}

	resp := s.PostJson(t, "/", data)
	resp.Ok(t)

	require.Equal(t, data.Nya, resp.Body().String())
}

func TestPostString(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		w.Write(buf.Bytes())
	})

	s := server.NewMockServer(mux)
	defer s.Close()

	body := "hello"

	resp := s.PostString(t, "/", "plain/text", body)
	resp.Ok(t)

	require.Equal(t, body, resp.Body().String())
}

func TestPostFormData(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		m := contents.NewMultipart()
		err := m.Insert("key", "value")
		require.NoError(t, err)

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			v := r.PostFormValue("key")

			w.Write([]byte(v))
		})

		s := server.NewMockServer(mux)
		defer s.Close()

		resp := s.PostFormData(t, "/", m)
		resp.Ok(t)

		require.Equal(t, "value", resp.Body().String())
	})

	t.Run("file", func(t *testing.T) {
		file, err := os.Open("../../README.md")
		require.NoError(t, err)

		m := contents.NewMultipart()
		err = m.InsertFile("file", file)
		require.NoError(t, err)

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			formFile, _, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(400)
				return
			}
			err = formFile.Close()
			if err != nil {
				w.WriteHeader(400)
				return
			}
		})

		s := server.NewMockServer(mux)
		defer s.Close()

		resp := s.PostFormData(t, "/", m)
		resp.Ok(t)
	})
}

func TestFormData(t *testing.T) {
	m := contents.NewMultipart()
	err := m.Insert("key", "value")
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(400)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		v := r.PostFormValue("key")

		w.Write([]byte(v))
	})

	s := server.NewMockServer(mux)
	defer s.Close()

	resp := s.FormData(t, "/", http.MethodPut, m)
	resp.Ok(t)

	require.Equal(t, "value", resp.Body().String())
}

func TestDo(t *testing.T) {
	t.Run("DELETE", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				w.WriteHeader(400)
				return
			}
		})

		s := server.NewMockServer(mux)
		defer s.Close()

		resp := s.Do(t, "/", http.MethodDelete, nil, func(r *http.Request) {})
		resp.Ok(t)
	})
}
