package config

import "github.com/spf13/viper"

func initDefaults() {
	viper.SetDefault("api.hostname", "")
	viper.SetDefault("api.port", "")
	viper.SetDefault("api.secret", "")

	viper.SetDefault("oauth.google.id", "")
	viper.SetDefault("oauth.google.secret", "")

	viper.SetDefault("discord.token", "")
}
