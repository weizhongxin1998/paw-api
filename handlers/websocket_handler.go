package handlers

import (
	"context"
	"encoding/json"
	"paw-api/pkg/httpclient"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WebSocketHandler struct {
	ctx context.Context
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

func (h *WebSocketHandler) SetContext(ctx context.Context) {
	h.ctx = ctx
}

func (h *WebSocketHandler) Connect(url string) error {
	return httpclient.WSConnect(url, func(msg httpclient.WSMessage) {
		if h.ctx == nil {
			return
		}
		data, _ := json.Marshal(msg)
		runtime.EventsEmit(h.ctx, "ws-message", string(data))
	})
}

func (h *WebSocketHandler) Send(content string) error {
	return httpclient.WSSend(content)
}

func (h *WebSocketHandler) Disconnect() {
	httpclient.WSClose()
}

func (h *WebSocketHandler) IsConnected() bool {
	return httpclient.WSIsConnected()
}
