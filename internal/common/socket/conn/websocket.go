package conn

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket webSocket连接
type WebSocket struct {
	msgType int
	conn    *websocket.Conn
}

var _ Conn = (*WebSocket)(nil)

// NewWebSocket 创建WebSocket
func NewWebSocket(conn *websocket.Conn, msgType int) *WebSocket {
	w := &WebSocket{
		conn:    conn,
		msgType: msgType,
	}

	return w
}

// ReadMessage 读取消息
func (w *WebSocket) ReadMessage() (p []byte, err error) {
	_, p, err = w.conn.ReadMessage()

	return
}

// Write 写入消息
func (w *WebSocket) Write(p []byte) (n int, err error) {
	err = w.conn.WriteMessage(w.msgType, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close 关闭连接
func (w *WebSocket) Close() error {
	return w.conn.Close()
}

// SetWriteDeadline 设置写入超时时间
func (w *WebSocket) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}

// SetReadDeadline 设置读取超时时间
func (w *WebSocket) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

// SetReadLimit 设置读取限制
func (w *WebSocket) SetReadLimit(limit int64) {
	w.conn.SetReadLimit(limit)
}

// LocalAddr 本机地址
func (w *WebSocket) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

// RemoteAddr 远程地址
func (w *WebSocket) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}
