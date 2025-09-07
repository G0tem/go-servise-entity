package factory

import (
	"fmt"

	"github.com/G0tem/go-servise-entity/internal/config"
	"gorm.io/gorm"
)

func startWebSocketClient(cfg *config.Config, db *gorm.DB) {
	fmt.Println("startWebSocketClient")
	// Создаём клиент один раз
	// client := payment.NewWebSocketClient(cfg.ApiToken, db)

	// // Запускаем подключение в фоне
	// go func() {
	// 	for {
	// 		// Пытаемся подключиться
	// 		if err := client.Connect(context.Background()); err != nil {
	// 			log.Printf("WebSocket connection error: %v. Retrying in 5 seconds...", err)
	// 			time.Sleep(5 * time.Second)
	// 			continue
	// 		}

	// 		// После успешного подключения ждём, пока соединение не оборвётся
	// 		<-client.Done()

	// 		// Как только соединение потеряно — начинаем цикл заново
	// 		log.Println("WebSocket disconnected. Reconnecting...")
	// 	}
	// }()
}
