package utils

import (
	"tg_shop/internal/service"
	"time"
)

func StartEarningProcessor(earningService service.Earning) {
	go func() {
		for {
			if err := earningService.ProcessEarnings(); err != nil {
				println("Error processing earnings:", err.Error())
			}

			time.Sleep(30 * time.Minute)
		}
	}()
}

func StartCheckPremiums(earningService service.Premium) {
	go func() {
		for {
			if err := earningService.GetPremiumInfo(); err != nil {
				println("Error processing earnings:", err.Error())
			}

			time.Sleep(30 * time.Minute)
		}
	}()
}
