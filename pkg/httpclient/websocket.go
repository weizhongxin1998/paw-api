package httpclient

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string `json:"type"`    // "sent" or "received"
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type WSConn struct {
	conn  *websocket.Conn
	mu    sync.Mutex
	done  chan struct{}
	onMsg func(WSMessage)
}

var ws *WSConn

func WSConnect(url string, onMessage func(WSMessage)) error {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	if ws != nil {
		WSClose()
	}
	ws = &WSConn{
		conn:  c,
		done:  make(chan struct{}),
		onMsg: onMessage,
	}
	go ws.readLoop()
	return nil
}

func (w *WSConn) readLoop() {
	for {
		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			close(w.done)
			return
		}
		if w.onMsg != nil {
			w.onMsg(WSMessage{
				Type:    "received",
				Content: string(msg),
				Time:    time.Now().UnixMilli(),
			})
		}
	}
}

func WSSend(content string) error {
	if ws == nil {
		return nil
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		return err
	}
	if ws.onMsg != nil {
		ws.onMsg(WSMessage{
			Type:    "sent",
			Content: content,
			Time:    time.Now().UnixMilli(),
		})
	}
	return nil
}

func WSClose() {
	if ws == nil {
		return
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.conn.Close()
	ws = nil
}

func WSIsConnected() bool {
	return ws != nil
}
