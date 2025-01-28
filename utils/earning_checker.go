package utils

import (
	"tg_shop/internal/service"
	"time"
)

func StartEarningProcessor(earningService service.Earning) {
	go func() {
		for {
			if err := earningService.ProcessEarnings(); err != nil {
				// Логируем ошибку (можно использовать log.Println или сторонний логгер)
				println("Error processing earnings:", err.Error())
			}

			time.Sleep(30 * time.Minute) // Пауза между обработкой партий
		}
	}()
}
