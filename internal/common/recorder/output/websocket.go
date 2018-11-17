package output

import (
	"github.com/ouqiang/mars/internal/common/recorder"
	"github.com/ouqiang/mars/internal/common/recorder/output/action"
	"github.com/ouqiang/mars/internal/common/socket"
	"github.com/ouqiang/mars/internal/common/socket/codec"
	"github.com/ouqiang/mars/internal/common/socket/message"
	log "github.com/sirupsen/logrus"
)

// WebSocket 输出到WebSocket
type WebSocket struct {
	router   *socket.Router
	hub      *socket.Hub
	codec    codec.Codec
	recorder *recorder.Recorder
}

// NewWebSocket 创建WebSocket
func NewWebSocket(hub *socket.Hub, r *recorder.Recorder) *WebSocket {
	ws := &WebSocket{
		router:   socket.NewRouter(),
		hub:      hub,
		codec:    new(codec.JSON),
		recorder: r,
	}

	ws.registerRouter()

	return ws
}

func (w *WebSocket) OnConnect(session *socket.Session) {
	log.Debugf("webSocket建立连接: [sessionId: %s]", session.ID)
	w.hub.Add(session.ID, session)
}
func (w *WebSocket) OnMessage(session *socket.Session, data []byte) {
	log.Debugf("webSocket收到消息: %s", data)
	w.router.Dispatch(session, w.codec, data)
}
func (w *WebSocket) OnClose(session *socket.Session) {
	log.Debugf("webSocket关闭连接: [sessionId: %s]", session.ID)
	w.hub.Delete(session.ID)
}
func (w *WebSocket) OnError(session *socket.Session, err error) {
	log.Debugf("webSocket错误: [sessionId: %s] %s", session.ID, err)
}

// Write Transaction写入WebSocket
func (w *WebSocket) Write(tx *recorder.Transaction) error {
	push := &action.PushTransaction{
		Id:       tx.Id,
		Method:   tx.Req.Method,
		Host:     tx.Req.Host,
		Path:     tx.Req.Path,
		Duration: tx.Duration,
	}
	if tx.Resp.Err != "" {
		push.ResponseErr = tx.Resp.Err
	} else {
		push.ResponseContentType = tx.Resp.Body.ContentType
		push.ResponseStatusCode = tx.Resp.StatusCode
		push.ResponseLen = tx.Resp.Body.Len
	}
	w.broadcast(action.TypePushTransaction, push)

	return nil
}

// 发送消息
func (w *WebSocket) sendMessage(session *socket.Session, msgType message.Type, payload interface{}) {
	data, err := w.marshalMessage(msgType, payload)
	if err != nil {
		return
	}
	session.Write(data)

}

// 广播
func (w *WebSocket) broadcast(msgType message.Type, payload interface{}) {
	data, err := w.marshalMessage(msgType, payload)
	if err != nil {
		return
	}
	w.hub.Broadcast(data)
}

// 消息序列化
func (w *WebSocket) marshalMessage(msgType message.Type, payload interface{}) ([]byte, error) {
	data, err := w.codec.Marshal(msgType, payload)
	if err != nil {
		log.Errorf("webSocket发送消息序列化错误: [msgType: %d payload: %+v] %s", msgType, payload, err)
		return nil, err
	}

	return data, nil
}

func (w *WebSocket) registerRouter() {
	w.router.NotFoundHandler(func(session *socket.Session, codec codec.Codec, data []byte) {
		log.Warnf("webSocket路由无法解析: [sessionId: %s] %s",
			session.ID, data)
	})
	w.router.ErrorHandler(func(session *socket.Session, err interface{}) {
		log.Warnf("webSocket路由解析错误: [sessionId: %s] %+v", session.ID, err)
	})

	w.router.Register(action.TypeRequestPing, (*action.Empty)(nil), w.ping)
	w.router.Register(action.TypeRequestReplay, (*action.RequestReplay)(nil), w.replay)
	w.router.Register(action.TypeRequestTransaction, (*action.RequestTransaction)(nil), w.getTransaction)
}

func (w *WebSocket) ping(ctx *socket.Context) {
	w.sendMessage(ctx.Session, action.TypeResponsePong, new(action.Empty))
}

func (w *WebSocket) replay(ctx *socket.Context) {
	req := ctx.Payload.(*action.RequestReplay)
	err := w.recorder.Replay(req.Id)
	resp := &action.ResponseReplay{}
	if err != nil {
		resp.Err = err.Error()
	}
	w.sendMessage(ctx.Session, action.TypeResponseReplay, resp)
}

func (w *WebSocket) getTransaction(ctx *socket.Context) {
	req := ctx.Payload.(*action.RequestTransaction)
	tx, err := w.recorder.Storage().Get(req.Id)
	resp := &action.ResponseTransaction{
		Transaction: tx,
	}
	if err != nil {
		resp.Err = err.Error()
	}
	w.sendMessage(ctx.Session, action.TypeResponseTransaction, resp)
}
