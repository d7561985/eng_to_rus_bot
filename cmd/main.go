package main

import (
	"crypto/tls"
	"github.com/d7561985/eng_to_rus_bot/pkg/multitran"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	go listener()
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Panic().Msg("BOT_TOKEN system variable not present")
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, &http.Client{Transport: &http.Transport{TLSClientConfig: getTLS()}})
	if err != nil {
		log.Panic().Err(err).Str("token", token).Msg("init bot api")
	}

	bot.Debug = os.Getenv("BOT_DEBUG") != ""

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		var msg tgbotapi.MessageConfig
		if len(update.Message.Text) > 50 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "too big request")
		} else {
			translate, err := multitran.GetWord(update.Message.Text)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			} else {

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, translate.String(120))
			}
		}

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Error().Err(err).Msgf("send message to: %s", update.Message.From.UserName)
		}
	}
}

func getTLS() *tls.Config {
	certificate, err := tls.LoadX509KeyPair("assets/cert.pem", "assets/key.pem")
	if err != nil {
		log.Panic().Err(err).Msg("load cetificates")
	}
	return &tls.Config{Certificates: []tls.Certificate{certificate}}
}

func listener() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	})
	if err := http.ListenAndServe(":"+port, http.DefaultServeMux); err != nil {
		log.Panic().Err(err).Msg("fake listener")
	}
}
