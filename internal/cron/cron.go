package cron

//
//import (
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/robfig/cron/v3"
//	"tg_shop/internal/repository"
//)
//
//type Cron struct {
//	PremiumsCron
//}
//
//func InitCron(bot *tgbotapi.BotAPI, repoUser repository.User) *cron.Cron {
//	cronScheduler := cron.New()
//	cronScheduler.Start()
//
//	premiumsCron := NewPremiumsCron(bot, repoUser)
//	cronScheduler.AddFunc("0 12 * * *", premiumsCron.CheckPremiumExpiry)
//
//	return cronScheduler
//}
//
//type Premiums interface {
//	CheckPremiumExpiry()
//}
