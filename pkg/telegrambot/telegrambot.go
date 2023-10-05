package telegrambot

import (
	"io"
	"log"
	"suntek2telegram/pkg/config"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramBot *tb.Bot

func Start(tc *config.Telegram, imgReadersChan chan io.Reader) {
	// Connect to Telegram bot
	var err error
	telegramBot, err = tb.NewBot(tb.Settings{
		Token:     tc.APIKey,
		Poller:    &tb.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tb.ModeMarkdownV2,
	})
	if err != nil {
		log.Fatalln("Failed to connect to Telegram bot:", err)
	}

	to := &tb.Chat{ID: int64(tc.ChatID)}
	log.Println("Telegram ready!")
	for reader := range imgReadersChan {
		fl := &tb.Photo{File: tb.FromReader(reader)}
		_, err := telegramBot.Send(to, fl)
		if err != nil {
			log.Fatalln("Failed to send to Telegram:", err)
		}
	}
}
