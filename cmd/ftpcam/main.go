package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"suntek2telegram/pkg/config"
	"suntek2telegram/pkg/ftpserver"
	"suntek2telegram/pkg/smtpserver"
	"suntek2telegram/pkg/telegrambot"
	"syscall"
)

var (
	flagConfigPath = flag.String("conf", "config.yml", "Path to config file")
)

func main() {
	flag.Parse()

	c, err := config.New(*flagConfigPath)
	if err != nil {
		log.Fatalln("Failed to load configuration file:", err)
	}

	imgReadersChan := make(chan io.Reader)

	if c.FTP.Enabled {
		go ftpserver.Start(c.FTP, imgReadersChan)
	}
	if c.SMTP.Enabled {
		go smtpserver.Start(c.SMTP, imgReadersChan)
	}
	go telegrambot.Start(c.Telegram, imgReadersChan)

	// Create a channel to wait for OS interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Block main function here until an interrupt is received
	<-interrupt
	fmt.Println("Program interrupted")
}
