package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"tg_shop/internal"
	"tg_shop/internal/handler"
	"tg_shop/internal/repository"
	"tg_shop/internal/service"
	"tg_shop/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Application started successfully")

	botToken := os.Getenv("BOT_TOKEN")

	bot := internal.InitBot(botToken)
	repos := repository.NewRepository(db)
	services := service.NewService(repos, bot)
	handlers := handler.NewHandler(services)
	adm_handlers := handler.NewAdminHandler(services)

	//cron.InitCron(bot, repos.User)

	go internal.BotProcess(handlers, bot)
	go internal.AdmBotProcess(adm_handlers)

	srv := new(internal.Server)
	if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running server %s", err.Error())
	}

	go utils.StartEarningProcessor(services.Earning)
	go utils.StartCheckPremiums(services.Premium)

	if err != nil {
		log.Panic(err)
	}
}
