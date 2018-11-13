package socket

import (
	"fmt"

	"github.com/ouqiang/mars/internal/common/socket/codec"
	"github.com/ouqiang/mars/internal/common/socket/message"
)

// Handler 路由处理器
type Handler func(ctx *Context)

// HandlersChain 路由处理器链
type HandlersChain []Handler

// router 路由
type Router struct {
	handlers        map[message.Type]HandlersChain
	handlersChain   HandlersChain
	registry        *message.Registry
	errorHandler    func(session *Session, err interface{})
	notFoundHandler func(session *Session, codec codec.Codec, data []byte)
}

// NewRouter 创建实例
func NewRouter() *Router {
	r := &Router{
		handlers: make(map[message.Type]HandlersChain),
		registry: message.NewRegistry(),
	}

	return r
}

// Register 路由注册
func (r *Router) Register(msgType message.Type, payload interface{}, handler Handler) *Router {
	if payload == nil {
		panic("payload type is nil")
	}
	if handler == nil {
		panic("handler is nil")
	}
	mergeHandlers := make(HandlersChain, len(r.handlersChain), len(r.handlersChain)+1)
	copy(mergeHandlers, r.handlersChain)
	mergeHandlers = append(mergeHandlers, handler)
	r.handlers[msgType] = mergeHandlers
	r.registry.Add(msgType, payload)

	return r
}

// Use 设置全局中间件
func (r *Router) Use(middleware ...Handler) {
	r.handlersChain = append(r.handlersChain, middleware...)
}

// Dispatch 路由分发
func (r *Router) Dispatch(session *Session, codec codec.Codec, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			r.triggerError(session, err)
		}
	}()
	msgType, payload, err := codec.Unmarshal(data, r.registry)
	if err != nil {
		if r.notFoundHandler != nil {
			r.notFoundHandler(session, codec, data)
			return
		}
		r.triggerError(session, fmt.Errorf("路由分发, 解码错误: %s", err))
		return
	}
	handlersChain, ok := r.handlers[msgType]
	if !ok {
		r.triggerError(session, fmt.Errorf("路由分发, 根据消息类型找不到对应的handler 消息类型: [%d] payload: [%+v]", msgType, payload))
		return
	}
	ctx := newContext()
	ctx.handlersChain = handlersChain
	ctx.Session = session
	ctx.MsgType = msgType
	ctx.Payload = payload
	ctx.Next()
}

// ErrorHandler 错误处理
func (r *Router) ErrorHandler(h func(session *Session, err interface{})) {
	r.errorHandler = h
}

// NotFoundHandler 未找到handler
func (r *Router) NotFoundHandler(h func(session *Session, codec codec.Codec, data []byte)) {
	r.notFoundHandler = h
}

func (r *Router) triggerError(session *Session, err interface{}) {
	if r.errorHandler != nil {
		r.errorHandler(session, err)
	}
}
