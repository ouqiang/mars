package controller

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Inspector 流量审查
type Inspector struct {
	Controller
}

// NewInspector 创建Inspector
func NewInspector() *Inspector {
	c := &Inspector{}

	return c
}

// WebSocket 处理webSocket
func (c *Inspector) WebSocket(resp http.ResponseWriter, req *http.Request) {
	upgrader.Upgrade(resp, req, nil)
}
