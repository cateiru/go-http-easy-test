package contents

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Multipart struct {
	body   *bytes.Buffer
	writer *multipart.Writer
}

func NewMultipart() *Multipart {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	return &Multipart{
		body:   body,
		writer: writer,
	}
}

func (c *Multipart) Insert(t *testing.T, key string, value string) {
	b := strings.NewReader(value)

	part, err := c.writer.CreateFormField(key)
	require.NoError(t, err)

	_, err = io.Copy(part, b)
	require.NoError(t, err)
}

func (c *Multipart) InsertFile(t *testing.T, key string, file *os.File) {
	part, err := c.writer.CreateFormFile(key, filepath.Base(file.Name()))
	require.NoError(t, err)

	_, err = io.Copy(part, file)
	require.NoError(t, err)
}

func (c *Multipart) Export() *bytes.Buffer {
	c.writer.Close()

	return c.body
}

func (c *Multipart) ContentType() string {
	return c.writer.FormDataContentType()
}
