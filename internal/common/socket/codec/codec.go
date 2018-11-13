package codec

import (
	"github.com/ouqiang/mars/internal/common/socket/message"
)

// Codec 编解码器
type Codec interface {
	// Marshal 编码
	Marshal(msgType message.Type, payload interface{}) ([]byte, error)
	// Unmarshal 解码
	Unmarshal(data []byte, registry *message.Registry) (msgType message.Type, payload interface{}, err error)
	// String 名称
	String() string
}
