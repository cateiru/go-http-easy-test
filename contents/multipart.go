package contents

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

func (c *Multipart) Insert(key string, value string) error {
	b := strings.NewReader(value)

	part, err := c.writer.CreateFormField(key)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, b)
	if err != nil {
		return err
	}

	return nil
}

func (c *Multipart) InsertFile(key string, file *os.File) error {
	part, err := c.writer.CreateFormFile(key, filepath.Base(file.Name()))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

func (c *Multipart) Export() *bytes.Buffer {
	c.writer.Close()

	return c.body
}

func (c *Multipart) ContentType() string {
	return c.writer.FormDataContentType()
}
