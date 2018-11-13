// Package inspector 流量审查
package inspector

import (
	"log"
	"net/http"

	"github.com/rakyll/statik/fs"

	"github.com/ouqiang/mars/internal/app/inject"
	"github.com/ouqiang/mars/internal/app/inspector/controller"
	_ "github.com/ouqiang/mars/internal/statik"
)

const staticDir = "/public/"

// router 路由
type Router struct {
	container *inject.Container
	mux       *http.ServeMux
}

// NewRouter 创建Router
func NewRouter(container *inject.Container, mux *http.ServeMux) *Router {
	r := &Router{
		container: container,
		mux:       mux,
	}

	return r
}

// Register 路由注册
func (r *Router) Register() {
	r.registerStatic()
	c := controller.NewInspector(r.container.WebSocketOutput, r.container.WebSocketSessionOpts)

	r.mux.HandleFunc("/ws", c.WebSocket)
}

func (r *Router) registerStatic() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r.mux.Handle(staticDir, http.StripPrefix(staticDir, http.FileServer(statikFS)))
}
