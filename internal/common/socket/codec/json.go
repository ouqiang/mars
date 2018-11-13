package codec

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ouqiang/mars/internal/common/socket/message"
)

// JSON json格式
type JSON struct {
	// Type 类型
	Type message.Type `json:"type"`
	// Payload 消息体
	Payload json.RawMessage `json:"payload"`
}

// Marshal 编码
func (JSON) Marshal(msgType message.Type, payload interface{}) ([]byte, error) {
	v, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	msg := &JSON{
		Type:    msgType,
		Payload: v,
	}

	return json.Marshal(msg)
}

// Unmarshal 解码
func (JSON) Unmarshal(data []byte, registry *message.Registry) (msgType message.Type, payload interface{}, e error) {
	if registry == nil {
		return 0, nil, errors.New("registry is nil")
	}
	jsonMsg := &JSON{}
	err := json.Unmarshal(data, jsonMsg)
	if err != nil {
		return 0, nil, err
	}
	value, found := registry.New(jsonMsg.Type)
	if !found {
		return 0, nil, fmt.Errorf("message type [%d] not found from registry", jsonMsg.Type)
	}
	err = json.Unmarshal(jsonMsg.Payload, value)
	if err != nil {
		return 0, nil, err
	}

	return jsonMsg.Type, value, err
}

func (JSON) String() string {
	return "json"
}
