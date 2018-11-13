package socket

import (
	"testing"

	"github.com/ouqiang/mars/internal/common/socket/codec"
	"github.com/ouqiang/mars/internal/common/socket/message"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	var msgType message.Type = 1000
	type Foo struct {
		Name string
	}
	j := &codec.JSON{}
	foo := &Foo{
		Name: "router",
	}
	data, err := j.Marshal(msgType, foo)
	require.NoError(t, err)
	r := NewRouter()

	sess := NewSession(nil)

	r.Use(func(ctx *Context) {
		ctx.Session.Data.Store("beforeNum", 1)
	}, func(ctx *Context) {
		ctx.Session.Data.Store("beforeNum", 2)
	}, func(ctx *Context) {
		ctx.Session.Data.Store("beforeNum", 3)
	})

	r.Register(msgType, (*Foo)(nil), func(ctx *Context) {
		beforeNum, _ := ctx.Session.Data.Load("beforeNum")
		require.Equal(t, 3, beforeNum)
		require.Equal(t, sess, ctx.Session)
		require.Equal(t, msgType, ctx.MsgType)
		f, ok := ctx.Payload.(*Foo)
		require.True(t, ok)
		require.Equal(t, foo, f)
	})

	r.Dispatch(sess, j, data)
	r.Register(msgType, (*Foo)(nil), func(ctx *Context) {
		panic("panic trigger")
	})
	errorNum := 0
	r.ErrorHandler(func(session *Session, payload interface{}) {
		require.Equal(t, sess, session)
		require.NotNil(t, payload)
		errorNum = 1
	})
	r.Dispatch(sess, j, data)
	require.Equal(t, 1, errorNum)
}
