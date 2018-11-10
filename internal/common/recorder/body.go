package recorder

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Body HTTP请求、响应Body
type Body struct {
	IsBinary    bool   `json:"is_binary"`
	Len         int    `json:"len"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}

// NewBody 创建Body
func NewBody() *Body {
	b := &Body{}

	return b
}

// 设置body内容
func (b *Body) setContent(contentType string, content []byte) {
	b.IsBinary = IsBinaryBody(contentType)
	b.Len = len(content)
	b.Content = content
	b.ContentType = contentType
}

// body内容封装成ReadCloser
func (b *Body) readCloser() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b.Content))
}
