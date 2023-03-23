package easy

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type Multipart struct {
	body   *bytes.Buffer
	writer *multipart.Writer
}

// Create a new multipart/form-data object
//
// Example:
//
//	m := NewMultipart()
//	// Insert k-v data
//	err := m.Insert("key", "value")
//	require.NoError(t, err)
//	// Insert files
//	err := m.InsertFile("key", file)
func NewMultipart() *Multipart {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	return &Multipart{
		body:   body,
		writer: writer,
	}
}

// Add a string form
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

// Add a file objects
//
// Example:
//
//	file, err := os.Open("file path")
//	require.NoError(t, err)
//	m.InsertFile("file", file)
func (c *Multipart) InsertFile(key string, file *os.File) error {
	contentType := detectFileContentType(file)
	fileName := filepath.Base(file.Name())

	mh := make(textproto.MIMEHeader)
	mh.Set("Content-Type", contentType)
	mh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, key, fileName))

	part, err := c.writer.CreatePart(mh)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

// Outputs a multipart/form-data format.
func (c *Multipart) Export() *bytes.Buffer {
	c.writer.Close()

	return c.body
}

// Outputs content-type
//
// ref. https://www.microfocus.com/documentation/idol/IDOL_12_0/MediaServer/Guides/html/English/Content/Shared_Admin/_ADM_POST_requests.htm#:~:text=In%20the%20multipart%2Fform%2Ddata,the%20data%20in%20the%20part.
func (c *Multipart) ContentType() string {
	return c.writer.FormDataContentType()
}

func detectFileContentType(file *os.File) string {
	buffer := make([]byte, 512)
	file.Read(buffer)

	contentType := http.DetectContentType(buffer)

	file.Seek(0, 0)

	return contentType
}
