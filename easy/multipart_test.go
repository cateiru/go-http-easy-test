package easy_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/cateiru/go-http-easy-test/v2/easy"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	t.Run("one element", func(t *testing.T) {
		m := easy.NewMultipart()

		err := m.Insert("key", "value")
		require.NoError(t, err)
	})

	t.Run("multi elements", func(t *testing.T) {
		m := easy.NewMultipart()

		err := m.Insert("key1", "value1")
		require.NoError(t, err)
		err = m.Insert("key2", "value2")
		require.NoError(t, err)
	})

	t.Run("same keys", func(t *testing.T) {
		m := easy.NewMultipart()

		err := m.Insert("key", "value1")
		require.NoError(t, err)
		err = m.Insert("key", "value2")
		require.NoError(t, err)
	})
}

func TestMultipart(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		m := easy.NewMultipart()

		data := map[string]string{
			"key":  "value",
			"111":  "aaaa",
			"mail": "test@example.com",
			"jp":   "日本語",
		}

		for key, value := range data {
			err := m.Insert(key, value)
			require.NoError(t, err)
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

		m := easy.NewMultipart()
		err = m.InsertFile("file", file)
		require.NoError(t, err)

		formData := m.Export().String()

		rep := regexp.MustCompile(`Content-Disposition: form-data; name="([^"]+)"; filename="([^"]+)"\r?\nContent-Type: ([^\r^\n]+)`)
		result := rep.FindAllStringSubmatch(formData, -1)

		require.Len(t, result, 1)

		require.Equal(t, result[0][1], "file")
		require.Equal(t, result[0][2], "README.md")
		require.Equal(t, result[0][3], "application/octet-stream")
	})

	t.Run("content-type", func(t *testing.T) {
		m := easy.NewMultipart()

		err := m.Insert("key", "value")
		require.NoError(t, err)

		r := regexp.MustCompile(`multipart/form-data; boundary=.+`)

		require.True(t, r.Match([]byte(m.ContentType())))
	})
}
