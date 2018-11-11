package output

import "github.com/ouqiang/mars/internal/common/recorder"

// WebSocket 输出到WebSocket
type WebSocket struct {
}

// Write Transaction写入WebSocket
func (w *WebSocket) Write(tx *recorder.Transaction) error {
	return nil
}
