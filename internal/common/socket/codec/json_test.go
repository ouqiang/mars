package codec

import (
	"testing"

	"github.com/ouqiang/mars/internal/common/socket/message"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	type Foo struct {
		Name string `json:"name"`
	}
	var msgType message.Type = 1
	j := JSON{}

	data, err := j.Marshal(msgType, &Foo{
		Name: "json",
	})
	require.NoError(t, err)
	require.Equal(t, []byte(`{"type":1,"payload":{"name":"json"}}`), data)
	registry := message.NewRegistry()
	newMsgType, payload, err := j.Unmarshal(data, registry)
	require.Error(t, err)

	registry.Add(msgType, (*Foo)(nil))
	newMsgType, payload, err = j.Unmarshal(data, registry)
	require.NoError(t, err)
	require.Equal(t, msgType, newMsgType)
	foo, ok := payload.(*Foo)
	require.True(t, ok)
	require.Equal(t, "json", foo.Name)
}
