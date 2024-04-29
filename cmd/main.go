package main

import (
	"bot/internal/config"
	"bot/internal/telegram"
	"log"
)

func main() {
	cfg, err := config.ReadConfigFromFile("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	bot := telegram.NewBot(cfg)
	bot.Start()
}
