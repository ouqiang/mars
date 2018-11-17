package action

import (
	"time"

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
	Err string `json:"err"`
}

type RequestTransaction struct {
	Id string `json:"id"`
}

type ResponseTransaction struct {
	*recorder.Transaction
	Err string `json:"err"`
}

type PushTransaction struct {
	Id string `json:"id"`
	// Method 请求方法
	Method string `json:"method"`
	// Host 请求主机名
	Host string `json:"host"`
	// Path 请求path
	Path string `json:"path"`
	// Duration 耗时
	Duration time.Duration `json:"duration"`
	// ResponseStatusCode 响应状态码
	ResponseStatusCode int `json:"response_status_code"`
	// Err 错误信息
	ResponseErr string `json:"response_err"`
	// ResponseContentType 响应内容类型
	ResponseContentType string `json:"response_content_type"`
	// ResponseLen 响应长度
	ResponseLen int `json:"response_len"`
}
