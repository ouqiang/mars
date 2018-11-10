package recorder

import (
	"net/http"

	"github.com/ouqiang/goproxy"
)

// Request HTTP请求
type Request struct {
	// Proto HTTP协议版本
	Proto string `json:"proto"`
	// Method 请求方法
	Method string `json:"method"`
	// Scheme 请求协议
	Scheme string `json:"scheme"`
	// Host 请求主机名
	Host string `json:"host"`
	// Path 请求path
	Path string `json:"path"`
	// QueryParam URL参数
	QueryParam string `json:"query_param"`
	// URL 完整URL
	URL string `json:"url"`
	// Header 请求Header
	Header http.Header `json:"header"`
	// Body 请求body
	Body *Body `json:"body"`
}

// NewRequest 创建请求
func NewRequest() *Request {
	req := &Request{
		Body: NewBody(),
	}

	return req
}

// Restore 还原请求
func (req *Request) Restore() (*http.Request, error) {
	rawReq, err := http.NewRequest(req.Method, req.URL, req.Body.readCloser())
	if err != nil {
		return nil, err
	}
	rawReq.Header = goproxy.CloneHeader(req.Header)

	return rawReq, nil
}
