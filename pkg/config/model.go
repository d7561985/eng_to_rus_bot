package config

// M is config model
type M struct {
	Port       string // Heroku default port env.
	HerokuSlug string // HEROKUSLUG => HEROKU_APP_NAME
	BotToken   string // BOTTOKEN => BOT_TOKEN
	BotDebug   bool   // BOTDEBUG => BOT_DEBUG
}
