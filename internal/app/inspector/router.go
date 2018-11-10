// Package inspector 流量审查
package inspector

import (
	"net/http"

	"github.com/ouqiang/mars/internal/app/inspector/controller"

	"github.com/ouqiang/mars/internal/app/inject"
)

// Router 路由
type Router struct {
	container *inject.Container
}

// NewRouter 创建路由管理器
func NewRouter(container *inject.Container) *Router {
	r := &Router{
		container: container,
	}

	return r
}

// Register 路由注册
func (r *Router) Register(mux *http.ServeMux) {
	c := controller.NewInspector(r.container.TxStorage)

	mux.HandleFunc("/ws", c.WebSocket)
	mux.HandleFunc("/get-transaction", c.GetTransaction)
}
