package socket

import (
	"github.com/ouqiang/mars/internal/common/socket/message"
)

const (
	abortIndex = 256
)

// Context 上下文信息
type Context struct {
	// Session 对应客户端连接
	Session *Session
	// MsgType 消息类型
	MsgType message.Type
	// Payload 消息体
	Payload       interface{}
	handlersChain HandlersChain
	index         int
}

func newContext() *Context {
	ctx := &Context{
		index: -1,
	}

	return ctx
}

// Abort 中断执行
func (ctx *Context) Abort() {
	ctx.index = abortIndex
}

// IsAborted 是否中断
func (ctx *Context) IsAborted() bool {
	return ctx.index >= abortIndex
}

// Next 执行下一个中间件
func (ctx *Context) Next() {
	ctx.index++
	l := len(ctx.handlersChain)
	for ; ctx.index < l; ctx.index++ {
		ctx.handlersChain[ctx.index](ctx)
	}
}
