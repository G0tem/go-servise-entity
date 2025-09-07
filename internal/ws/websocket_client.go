package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/G0tem/go-servise-entity/internal/model"
	"github.com/G0tem/go-servise-entity/internal/types"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	wsURLTemplate = "wss://ws-gate.invoicebox.ru/v3/gate/rpc?token=%s"
	writeWait     = 10 * time.Second
	pongWait      = 60 * time.Second
	pingPeriod    = (pongWait * 9) / 10
)

type InvoiceboxWebSocketClient struct {
	conn      *websocket.Conn
	apiToken  string
	db        *gorm.DB
	done      chan struct{}
	interrupt chan struct{}
}

func NewInvoiceboxWebSocketClient(apiToken string, db *gorm.DB) *InvoiceboxWebSocketClient {
	return &InvoiceboxWebSocketClient{
		apiToken:  apiToken,
		db:        db,
		done:      make(chan struct{}),
		interrupt: make(chan struct{}),
	}
}

func (c *InvoiceboxWebSocketClient) Connect(ctx context.Context) error {
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

func (c *InvoiceboxWebSocketClient) Done() <-chan struct{} {
	return c.done
}

func (c *InvoiceboxWebSocketClient) close() {
	close(c.done)
}

func (c *InvoiceboxWebSocketClient) readPump() {
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
			Method string                  `json:"method"`
			Params types.OrderNotification `json:"params"`
		}

		if err := json.Unmarshal(message, &notification); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if notification.Method == "onOrderStatusChange" {
			log.Printf("Order status changed: ID=%s, Status=%s", notification.Params.ID, notification.Params.Status)

			// Обновляем статус в БД
			err := updatePaymentStatus(c.db, notification.Params)
			if err != nil {
				log.Printf("Failed to update payment status: %v", err)
			}
		} else if notification.Method == "getStatus" {
			// Отвечаем диагностикой
			resp := map[string]interface{}{
				"method": "status",
				"result": "OK",
			}
			_ = c.conn.WriteJSON(resp)
		}
	}
}

func (c *InvoiceboxWebSocketClient) writePump() {
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

func updatePaymentStatus(db *gorm.DB, notif types.OrderNotification) error {
	var payment model.Entity
	if err := db.Where("invoicebox_id = ?", notif.ID).First(&payment).Error; err != nil {
		return fmt.Errorf("failed to find payment by invoicebox_id: %w", err)
	}

	// Обновляем статус
	return db.Model(&payment).Updates(map[string]interface{}{
		"status":      notif.Status,
		"payment_url": nil, // можно очистить, если заказ завершён
	}).Error
}
