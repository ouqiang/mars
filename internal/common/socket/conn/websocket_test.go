package conn

import (
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/require"
)

func testWebSocketHandler(msgType int) http.Handler {
	upgrade := websocket.Upgrader{}
	handler := func(rw http.ResponseWriter, req *http.Request) {
		conn, err := upgrade.Upgrade(rw, req, nil)
		if err != nil {
			panic(err)
		}
		webSocket := NewWebSocket(conn, msgType)
		data, err := webSocket.ReadMessage()
		if err != nil {
			panic(err)
		}
		_, err = webSocket.Write(data)
		if err != nil {
			panic(err)
		}
	}

	return http.HandlerFunc(handler)
}

func TestWebSocket(t *testing.T) {
	conn, resp, err := wstest.NewDialer(testWebSocketHandler(websocket.TextMessage)).Dial("ws:/", nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusSwitchingProtocols)
	msg := []byte("ping")
	err = conn.WriteMessage(websocket.TextMessage, msg)
	require.NoError(t, err)
	_, data, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, msg, data)
}
