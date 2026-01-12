package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	wsURLTemplate = "wss://ws-gate.service.ru/v3/gate/rpc?token=%s"
	writeWait     = 10 * time.Second
	pongWait      = 60 * time.Second
	pingPeriod    = (pongWait * 9) / 10
)

type WebSocketClient struct {
	conn      *websocket.Conn
	apiToken  string
	db        *gorm.DB
	done      chan struct{}
	interrupt chan struct{}
}

func NewWebSocketClient(apiToken string, db *gorm.DB) *WebSocketClient {
	return &WebSocketClient{
		apiToken:  apiToken,
		db:        db,
		done:      make(chan struct{}),
		interrupt: make(chan struct{}),
	}
}

func (c *WebSocketClient) Connect(ctx context.Context) error {
	url := fmt.Sprintf(wsURLTemplate, c.apiToken)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.readPump()
	go c.writePump()

	log.Printf("Success connect to WS system InvoiceBox")

	return nil
}

func (c *WebSocketClient) Done() <-chan struct{} {
	return c.done
}

func (c *WebSocketClient) readPump() {
	// defer сработает, когда readPump завершится (при разрыве соединения)
	defer func() {
		close(c.done)      // уведомляем, что клиент отключён
		_ = c.conn.Close() // закрываем соединение
	}()

	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		log.Printf("Received raw message: %s", message)

		var notification struct {
			Method string `json:"method"`
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
	}
}

func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.interrupt:
			return
		}
	}
}
