package main

import (
	"context"
	"log"
	"os"
	"teko_pracrice_bot/bot"
	"teko_pracrice_bot/store"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}
	webHookURL := os.Getenv("WEBHOOK_URL")
	if botToken == "" {
		log.Fatal("$WEBHOOK_URL must be set")
	}

	storage := store.New(databaseUrl)
	err = storage.Open()
	if err != nil {
		log.Panic(err)
	}
	defer storage.Close()

	botAPI, err := bot.GetAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot := bot.New(storage, botAPI)

	err = bot.SetWebHook(webHookURL)
	if err != nil {
		log.Panic(err)
	}

	err = bot.Start(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
