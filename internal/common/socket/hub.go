package socket

import (
	"io"
	"sync"
	"sync/atomic"
)

// Hub session管理
type Hub struct {
	sessions  sync.Map
	broadcast chan []byte
	num       int32
}

// NewHub 创建session集合实例
func NewHub(broadcastQueueSize int) *Hub {
	h := &Hub{
		broadcast: make(chan []byte, broadcastQueueSize),
	}

	go h.run()

	return h
}

// Get 获取session
func (ch *Hub) Get(sessionID string) (io.Writer, bool) {
	value, ok := ch.sessions.Load(sessionID)
	if !ok {
		return nil, false
	}

	return value.(io.Writer), true
}

// Add 添加session
func (ch *Hub) Add(sessionID string, w io.Writer) {
	if sessionID == "" || w == nil {
		return
	}
	ch.sessions.Store(sessionID, w)
	atomic.AddInt32(&ch.num, 1)
}

// Delete 删除session
func (ch *Hub) Delete(sessionID string) {
	ch.sessions.Delete(sessionID)
	atomic.AddInt32(&ch.num, -1)
}

// Broadcast 广播
func (ch *Hub) Broadcast(data []byte) {
	ch.broadcast <- data
}

// Num session数量
func (ch *Hub) Num() int32 {
	return atomic.LoadInt32(&ch.num)
}

// Range 遍历session
func (ch *Hub) Range(f func(key, value interface{}) bool) {
	ch.sessions.Range(f)
}

// 运行
func (ch *Hub) run() {
	for data := range ch.broadcast {
		ch.sessions.Range(func(key, value interface{}) bool {
			w := value.(io.Writer)
			w.Write(data)
			return true
		})
	}
}
