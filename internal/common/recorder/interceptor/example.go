package interceptor

import (
	"io"
	"net/http"
	"strings"

	"github.com/ouqiang/goproxy"
)

func init() {
	//register(new(Example))
}

type Example struct{}

// Connect 收到客户端连接, 自定义response返回
// NOTICE: HTTPS只能访问 ctx.Req.URL.Host, 不能访问Header和Body, 不能使用rw
func (Example) Connect(ctx *goproxy.Context, rw http.ResponseWriter) {
	if strings.Contains(ctx.Req.URL.Host, "crashlytics.com") {
		rw.WriteHeader(http.StatusForbidden)
		ctx.Abort()
	}
	io.WriteString(rw, "修改后的内容")
	ctx.Abort()
}

// BeforeRequest 请求发送前, 修改request
func (Example) BeforeRequest(ctx *goproxy.Context) {
	ctx.Req.Header.Set("Req-Id", "123")
}

// BeforeResponse 响应发送前, 修改response
func (Example) BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error) {
	if err == nil {
		resp.Header.Set("Resp-Id", "456")
	}
}
