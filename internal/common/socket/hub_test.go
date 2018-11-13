package socket

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testSession struct {
	any int
}

func (tc *testSession) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestHub(t *testing.T) {
	hub := NewHub(50)

	var clientNum int64 = 10
	clients := make([]*testSession, 0, clientNum)
	var i int64 = 1
	for ; i <= clientNum; i++ {
		c := &testSession{}
		clients = append(clients, c)
		hub.Add(strconv.FormatInt(i, 10), c)
	}
	time.Sleep(1 * time.Millisecond)
	require.Equal(t, int32(clientNum), hub.Num())

	hub.Broadcast([]byte("session broadcast"))

	for i = 1; i <= clientNum; i++ {
		hub.Delete(strconv.FormatInt(i, 10))
	}
	time.Sleep(1 * time.Millisecond)
	require.Equal(t, int32(0), hub.Num())
}
