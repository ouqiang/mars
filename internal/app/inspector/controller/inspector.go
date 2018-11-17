package controller

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ouqiang/mars/internal/common/socket/conn"

	log "github.com/sirupsen/logrus"

	"github.com/ouqiang/mars/internal/common/socket"

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
	sessionOptions []socket.SessionOption
	sessionHandler socket.SessionHandler
}

// NewInspector 创建Inspector
func NewInspector(sessionHandler socket.SessionHandler, opts []socket.SessionOption) *Inspector {
	c := &Inspector{
		sessionHandler: sessionHandler,
		sessionOptions: opts,
	}

	return c
}

// WebSocket 处理webSocket
func (c *Inspector) WebSocket(resp http.ResponseWriter, req *http.Request) {
	rawConn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Debugf("升级到websocket错误: %s", err)
		return
	}
	client := socket.NewSession(
		conn.NewWebSocket(rawConn, websocket.TextMessage),
		c.sessionHandler,
		c.sessionOptions...,
	)
	client.ID = uuid.NewV4().String()
	client.Run()
}
