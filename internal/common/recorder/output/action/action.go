package action

import (
	"github.com/ouqiang/mars/internal/common/recorder"
	"github.com/ouqiang/mars/internal/common/socket/message"
)

const (
	TypeRequestPing        message.Type = 1000
	TypeRequestReplay      message.Type = 1001
	TypeRequestTransaction message.Type = 1002

	TypeResponsePong        message.Type = 2000
	TypeResponseReplay      message.Type = 2001
	TypeResponseTransaction message.Type = 2002

	TypePushTransaction message.Type = 3000
)

type Empty struct {
}

type RequestReplay struct {
	Id string `json:"id"`
}

type ResponseReplay struct {
	Err error `json:"err"`
}

type RequestTransaction struct {
	Id string `json:"id"`
}

type ResponseTransaction struct {
	*recorder.Transaction
	Err error `json:"err"`
}

type PushTransaction struct {
	Id string `json:"id"`
	// Method 请求方法
	Method string `json:"method"`
	// Scheme 请求协议
	Scheme string `json:"scheme"`
	// Host 请求主机名
	Host string `json:"host"`
	// Path 请求path
	Path string `json:"path"`
	// URL 完整URL
	URL string `json:"url"`
	// Status 状态码
	ResponseStatus string `json:"response_status"`
	// Err 错误信息
	ResponseErr error `json:"response_err"`
}
