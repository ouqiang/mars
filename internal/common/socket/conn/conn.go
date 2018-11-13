package conn

import (
	"io"
	"net"
	"time"
)

// Conn 连接
type Conn interface {
	// SetWriteDeadline 写入超时
	SetWriteDeadline(t time.Time) error

	// SetReadDeadline 读取超时
	SetReadDeadline(t time.Time) error

	// SetReadLimit 读取限制
	SetReadLimit(limit int64)

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远程地址
	RemoteAddr() net.Addr

	// ReadMessage 读取消息
	ReadMessage() (p []byte, err error)

	io.Writer
	io.Closer
}
