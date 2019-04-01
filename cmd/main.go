package main

import (
	"github.com/d7561985/eng_to_rus_bot/pkg/bot"
	"github.com/d7561985/eng_to_rus_bot/pkg/config"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	//go listener()

	bot.B().WebHook()
}

func listener() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//_, _ = writer.Write([]byte(strings.Join(os.Environ(), "\n")))
	})
	if err := http.ListenAndServe(":"+config.V.Port, http.DefaultServeMux); err != nil {
		log.Panic().Err(err).Msg("fake listener")
	}
}
