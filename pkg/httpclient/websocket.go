package httpclient

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn        *websocket.Conn
	isConnected bool
}

func NewWSClient() *WSClient {
	return &WSClient{}
}

func (w *WSClient) Connect(url string, headers map[string]string) error {
	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		return err
	}

	w.conn = conn
	w.isConnected = true
	return nil
}

func (w *WSClient) Disconnect() error {
	w.isConnected = false
	if w.conn != nil {
		err := w.conn.Close()
		w.conn = nil
		return err
	}
	return nil
}

func (w *WSClient) Send(message []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, message)
}

func (w *WSClient) Read() ([]byte, error) {
	_, msg, err := w.conn.ReadMessage()
	return msg, err
}

func (w *WSClient) IsConnected() bool {
	return w.isConnected
}
