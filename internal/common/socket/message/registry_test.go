package message

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	type Foo struct {
		Name string
	}

	var msgType Type = 10
	r.Add(msgType, (*Foo)(nil))
	v, found := r.New(msgType)
	require.True(t, found)
	require.Equal(t, reflect.TypeOf((*Foo)(nil)).Elem().Name(), reflect.TypeOf(v).Elem().Name())
	require.NotEqual(t, reflect.TypeOf((*Foo)(nil)).Elem().Name(), reflect.TypeOf((*Registry)(nil)).Elem().Name())

	v, found = r.New(Type(11))
	require.False(t, found)
	require.Nil(t, v)
}
