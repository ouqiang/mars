package recorder

import "net/http"

// Response HTTP响应
type Response struct {
	// Proto 响应协议
	Proto string `json:"proto"`
	// Status 状态码
	Status string `json:"status"`
	// Header 响应Header
	Header http.Header `json:"header"`
	// Body 响应Body
	Body *Body `json:"body"`
	// Err 错误信息
	Err error `json:"err"`
}

func NewResponse() *Response {
	resp := &Response{
		Body: NewBody(),
	}

	return resp
}
