package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ouqiang/mars/internal/common/recorder"
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
	storage recorder.Storage
}

// NewInspector 创建Inspector
func NewInspector(s recorder.Storage) *Inspector {
	c := &Inspector{
		storage: s,
	}

	return c
}

// WebSocket 处理webSocket
func (c *Inspector) WebSocket(resp http.ResponseWriter, req *http.Request) {
	upgrader.Upgrade(resp, req, nil)
}

// GetTransaction 获取Transaction
func (c *Inspector) GetTransaction(resp http.ResponseWriter, req *http.Request) {
	txId := req.FormValue("id")
	tx, err := c.storage.Get(txId)
	if err != nil {
		io.WriteString(resp, err.Error())
		return
	}
	data, err := json.Marshal(tx)
	if err != nil {
		io.WriteString(resp, err.Error())
		return
	}
	resp.Write(data)
}
