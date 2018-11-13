package codec

import (
	"testing"

	"github.com/ouqiang/mars/internal/common/socket/message"
	"github.com/stretchr/testify/require"
)

func TestProtobuf(t *testing.T) {
	var msgType message.Type = 100
	pb := &Foo{
		Name: "protobuf",
	}
	p := Protobuf{}
	data, err := p.Marshal(msgType, pb)
	require.NoError(t, err)

	registry := message.NewRegistry()
	newMsgType, payload, err := p.Unmarshal(data, registry)
	require.Error(t, err)

	registry.Add(msgType, (*Foo)(nil))
	newMsgType, payload, err = p.Unmarshal(data, registry)
	require.NoError(t, err)
	require.Equal(t, msgType, newMsgType)

	foo, ok := payload.(*Foo)
	require.True(t, ok)
	require.Equal(t, foo.GetName(), pb.Name)
}
