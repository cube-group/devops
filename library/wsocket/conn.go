package wsocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 1 * time.Hour

	ReadWait = 1 * time.Hour
	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
	// Maximum message size allowed from peer.
	MaxMessageSize = 512
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决websocket 403错误
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
