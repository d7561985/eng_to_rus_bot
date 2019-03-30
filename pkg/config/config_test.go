package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	_ = os.Setenv("BOT_TOKEN", "123")
	err := read()
	assert.NoError(t, err)
	assert.Equal(t, V.BotToken, "123")
	assert.Equal(t, viper.Get("BOT_TOKEN"), "123")
	assert.Equal(t, V.Port, "3000")
	assert.Equal(t, viper.Get("PORT"), "3000")
}
