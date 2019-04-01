package bot

import (
	"crypto/tls"
	"github.com/d7561985/eng_to_rus_bot/pkg/config"
	"github.com/d7561985/eng_to_rus_bot/pkg/multitran"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"net/http"
)

type bot struct {
}

func B() *bot {
	if config.V.BotToken == "" {
		log.Panic().Msg("BOT_TOKEN system variable not present")
	}
	return &bot{}
}

// LiveListener
func (b *bot) LiveListener() {
	bot, err := tgbotapi.NewBotAPIWithClient(config.V.BotToken, &http.Client{Transport: &http.Transport{TLSClientConfig: getTLS()}})
	if err != nil {
		log.Panic().Err(err).Str("token", config.V.BotToken).Msg("init bot api")
	}

	bot.Debug = config.V.BotDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal().Err(err).Msg("get updates")
	}

	log.Info().Msg("start")
	update(updates, bot)
}

// require meta data on heroku
// $ heroku labs:enable runtime-dyno-metadata -a {{.APP_NAME}}
// after this we can use HEROKU_APP_NAME for get slug domain in https://{{.HEROKU_APP_NAME}}.herokuapp.com/
// for registering
func (b *bot) WebHook() {
	if config.V.HerokuSlug == "" {
		log.Panic().Msg("HEROKU_APP_NAME not exist. There is no way to check domain name. Try to enable runtime-dyno-metadata")
	}

	bot, err := tgbotapi.NewBotAPI(config.V.BotToken)
	if err != nil {
		log.Panic().Err(err).Str("token", config.V.BotToken).Msg("init bot api")
	}

	bot.Debug = config.V.BotDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+config.V.HerokuSlug+".herokuapp.com/"+bot.Token, "assets/cert.pem"))
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://" + config.V.HerokuSlug + ".herokuapp.com/" + bot.Token))
	if err != nil {
		log.Fatal().Err(err).Msg("set hook")
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal().Err(err).Msg("get hook info")
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	log.Info().Msg("start")
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("0.0.0.0:"+config.V.Port, nil)

	update(updates, bot)
}

// update controller
func update(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
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
