package contents_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/cateiru/go-http-easy-test/contents"
	"github.com/stretchr/testify/require"
)

func TestMultipart(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		m := contents.NewMultipart()

		data := map[string]string{
			"key":  "value",
			"111":  "aaaa",
			"mail": "test@example.com",
			"jp":   "日本語",
		}

		for key, value := range data {
			m.Insert(t, key, value)
		}

		formData := m.Export().String()

		rep := regexp.MustCompile(`Content-Disposition: form-data; name="([^"]+)"\r?\n\r?\n([^\r^\n]+)\r?\n`)
		result := rep.FindAllStringSubmatch(formData, -1)

		require.Len(t, result, len(data))

		for _, d := range result {
			require.Equal(t, data[d[1]], d[2])
		}
	})

	t.Run("insert file", func(t *testing.T) {
		file, err := os.Open("../README.md")
		require.NoError(t, err)

		m := contents.NewMultipart()
		m.InsertFile(t, "file", file)

		formData := m.Export().String()

		t.Log(formData)

		rep := regexp.MustCompile(`Content-Disposition: form-data; name="([^"]+)"; filename="([^"]+)"\r?\nContent-Type: ([^\r^\n]+)`)
		result := rep.FindAllStringSubmatch(formData, -1)

		require.Len(t, result, 1)

		require.Equal(t, result[0][1], "file")
		require.Equal(t, result[0][2], "README.md")
		require.Equal(t, result[0][3], "application/octet-stream")
	})

	t.Run("content-type", func(t *testing.T) {
		m := contents.NewMultipart()

		m.Insert(t, "key", "value")

		r := regexp.MustCompile(`multipart/form-data; boundary=.+`)

		require.True(t, r.Match([]byte(m.ContentType())))
	})
}
