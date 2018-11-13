package socket

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ouqiang/mars/internal/common/socket/conn"
)

// 可选项
type options struct {
	// 读取最大字节数
	readLimit int64
	// 心跳超时时间
	heartBeatTimeout time.Duration
	// 接收队列大小
	receiveQueueSize int
	// 发送队列大小
	sendQueueSize int
	// 写入超时时间
	writeTimeout time.Duration
	// 读取超时时间
	readTimeout time.Duration
	// 心跳超时次数
	heartBeatTimeoutTimes int
}

var defaultSessionOptions = options{
	readLimit:             4 << 20,
	heartBeatTimeout:      30 * time.Second,
	receiveQueueSize:      50,
	sendQueueSize:         50,
	readTimeout:           2 * time.Minute,
	writeTimeout:          10 * time.Second,
	heartBeatTimeoutTimes: 2,
}

// SessionOption 可选项
type SessionOption func(*options)

// SessionHandler 事件处理
type SessionHandler interface {
	OnConnect(*Session)
	OnMessage(*Session, []byte)
	OnClose(*Session)
	OnError(*Session, error)
}

// WithSessionHeartBeatTimeoutTimes 心跳超时次数
func WithSessionHeartBeatTimeoutTimes(times int) SessionOption {
	return func(opt *options) {
		opt.heartBeatTimeoutTimes = times
	}
}

// WithSessionReadLimit 读取最大字节数
func WithSessionReadLimit(limit int64) SessionOption {
	return func(opt *options) {
		opt.readLimit = limit
	}
}

// WithSessionHeartBeatTimeout 心跳超时
func WithSessionHeartBeatTimeout(d time.Duration) SessionOption {
	return func(opt *options) {
		opt.heartBeatTimeout = d
	}
}

// WithSessionReceiveQueueSize 接收队列大小
func WithSessionReceiveQueueSize(size int) SessionOption {
	return func(opt *options) {
		opt.receiveQueueSize = size
	}
}

// WithSessionSendQueueSize 发送队列大小
func WithSessionSendQueueSize(size int) SessionOption {
	return func(opt *options) {
		opt.sendQueueSize = size
	}
}

// WithSessionWriteTimeout 写入超时
func WithSessionWriteTimeout(d time.Duration) SessionOption {
	return func(opt *options) {
		opt.writeTimeout = d
	}
}

// WithSessionReadTimeout 读取超时
func WithSessionReadTimeout(d time.Duration) SessionOption {
	return func(opt *options) {
		opt.readTimeout = d
	}
}

// Session 对应一个客户端连接
type Session struct {
	ID string
	// Data 保存Session相关数据
	Data           sync.Map
	handler        SessionHandler
	conn           conn.Conn
	lastActiveTime atomic.Value
	receiveChan    chan []byte
	sendChan       chan []byte
	closeChan      chan struct{}
	closeOnce      sync.Once
	closed         atomic.Value
	heartBeatTimer *time.Ticker
	opts           options
}

// NewSession 创建session实例
func NewSession(conn conn.Conn, handler SessionHandler, opt ...SessionOption) *Session {
	if conn == nil {
		panic("conn is nil")
	}
	if handler == nil {
		panic("handler is nil")
	}
	opts := defaultSessionOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &Session{
		handler:        handler,
		opts:           opts,
		conn:           conn,
		receiveChan:    make(chan []byte, opts.receiveQueueSize),
		sendChan:       make(chan []byte, opts.sendQueueSize),
		closeChan:      make(chan struct{}),
		heartBeatTimer: time.NewTicker(opts.heartBeatTimeout),
	}
	s.lastActiveTime.Store(time.Now())
	s.closed.Store(false)

	return s
}

// Run 运行
func (s *Session) Run() {
	go s.startHeartTimer()
	go s.receiveMessage()
	go s.handleMessage()
	go s.sendMessage()

	s.handler.OnConnect(s)
}

// Write 写入数据
func (s *Session) Write(data []byte) (n int, err error) {
	select {
	case s.sendChan <- data:
	default:
		s.Close()
	}

	return len(data), nil
}

// RemoteAddr 远程主机地址
func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

// Closed 连接是否已关闭
func (s *Session) Closed() bool {
	return s.closed.Load().(bool)
}

// CloseNotify 连接关闭通知
func (s *Session) CloseNotify() <-chan struct{} {
	return s.closeChan
}

// Close 关闭连接
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		close(s.closeChan)
		s.conn.Close()

		s.handler.OnClose(s)
	})
}

// 处理消息
func (s *Session) handleMessage() {
	for data := range s.receiveChan {
		s.handler.OnMessage(s, data)
	}
}

// 发送消息
func (s *Session) sendMessage() {
	defer s.Close()
	for {
		select {
		case <-s.closeChan:
			return
		case data := <-s.sendChan:
			s.conn.SetWriteDeadline(time.Now().Add(s.opts.writeTimeout))
			_, err := s.conn.Write(data)
			if err != nil {
				s.handler.OnError(s, err)
				return
			}
		}
	}
}

// 接收消息
func (s *Session) receiveMessage() {
	defer func() {
		close(s.receiveChan)
		s.Close()
	}()

	for {
		s.conn.SetReadDeadline(time.Now().Add(s.opts.readTimeout))
		s.conn.SetReadLimit(s.opts.readLimit)
		data, err := s.conn.ReadMessage()
		if err != nil {
			s.handler.OnError(s, err)
			return
		}
		s.lastActiveTime.Store(time.Now())
		s.receiveChan <- data
	}
}

// 启动心跳定时器
func (s *Session) startHeartTimer() {
	defer s.heartBeatTimer.Stop()

	timeoutTimes := 0
	for {
		select {
		case <-s.closeChan:
			return
		case <-s.heartBeatTimer.C:
			lastActiveTime := s.lastActiveTime.Load().(time.Time)
			if time.Now().Sub(lastActiveTime) <= s.opts.heartBeatTimeout {
				timeoutTimes = 0
				continue
			}
			timeoutTimes++
			if timeoutTimes >= s.opts.heartBeatTimeoutTimes {
				s.Close()
				return
			}
		}
	}
}
