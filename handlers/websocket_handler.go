package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"paw-api/pkg/httpclient"
)

type wsConnection struct {
	client *httpclient.WSClient
	ctx    context.Context
	url    string
}

type WebSocketHandler struct {
	ctx         context.Context
	connections sync.Map
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

func (h *WebSocketHandler) SetContext(ctx context.Context) {
	h.ctx = ctx
}

func (h *WebSocketHandler) Connect(url string, headersJSON string) error {
	if _, loaded := h.connections.Load(url); loaded {
		return fmt.Errorf("already connected to %s", url)
	}

	var headers map[string]string
	if headersJSON != "" {
		json.Unmarshal([]byte(headersJSON), &headers)
	}

	client := httpclient.NewWSClient()
	if err := client.Connect(url, headers); err != nil {
		return err
	}

	conn := &wsConnection{
		client: client,
		ctx:    h.ctx,
		url:    url,
	}

	h.connections.Store(url, conn)

	runtime.EventsEmit(h.ctx, "ws:connected", url)

	go h.readLoop(conn)

	return nil
}

func (h *WebSocketHandler) readLoop(conn *wsConnection) {
	for {
		msg, err := conn.client.Read()
		if err != nil {
			if conn.client.IsConnected() {
				runtime.EventsEmit(conn.ctx, "ws:error", conn.url, err.Error())
			}
			runtime.EventsEmit(conn.ctx, "ws:closed", conn.url)
			conn.client.Disconnect()
			h.connections.Delete(conn.url)
			return
		}

		runtime.EventsEmit(conn.ctx, "ws:message", conn.url, string(msg))
	}
}

func (h *WebSocketHandler) Send(url string, message string) error {
	val, ok := h.connections.Load(url)
	if !ok {
		return fmt.Errorf("not connected to %s", url)
	}
	conn := val.(*wsConnection)
	return conn.client.Send([]byte(message))
}

func (h *WebSocketHandler) Disconnect(url string) error {
	val, ok := h.connections.Load(url)
	if !ok {
		return fmt.Errorf("not connected to %s", url)
	}
	conn := val.(*wsConnection)
	conn.client.Disconnect()
	h.connections.Delete(url)
	runtime.EventsEmit(h.ctx, "ws:closed", url)
	return nil
}
