package config

import "github.com/spf13/viper"

func initDefaults() {
	viper.SetDefault("discord.token", "")
	viper.SetDefault("discord.app.id", "")

	viper.SetDefault("api.url", "")
	viper.SetDefault("api.port", "")
	viper.SetDefault("api.secret", "")

	viper.SetDefault("oauth.google.id", "")
	viper.SetDefault("oauth.google.secret", "")
}
