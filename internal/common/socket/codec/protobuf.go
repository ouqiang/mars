package codec

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ouqiang/mars/internal/common/socket/message"
)

// Protobuf protobuf格式
type Protobuf struct{}

// Marshal 编码
func (Protobuf) Marshal(msgType message.Type, payload interface{}) ([]byte, error) {
	pb, ok := payload.(proto.Message)
	if !ok {
		return nil, errors.New("invalid protobuf payload")
	}
	data, err := ptypes.MarshalAny(pb)
	if err != nil {
		return nil, err
	}
	msg := &message.Message{
		Type:    int32(msgType),
		Payload: data,
	}

	return proto.Marshal(msg)
}

// Unmarshal 解码
func (Protobuf) Unmarshal(data []byte, registry *message.Registry) (msgType message.Type, payload interface{}, e error) {
	if registry == nil {
		return 0, nil, errors.New("registry is nil")
	}
	msg := &message.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return 0, nil, err
	}
	v, found := registry.New(message.Type(msg.Type))
	if !found {
		return 0, nil, fmt.Errorf("message type [%d] not found from registry", msg.Type)
	}
	pb, ok := v.(proto.Message)
	if !ok {
		return 0, nil, fmt.Errorf("registry type [%d] is not pb", msg.Type)
	}
	err = ptypes.UnmarshalAny(msg.Payload, pb)
	if err != nil {
		return 0, nil, err
	}

	return message.Type(msg.Type), pb, err
}

func (Protobuf) String() string {
	return "protobuf"
}
