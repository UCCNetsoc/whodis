package config

import "github.com/spf13/viper"

func initDefaults() {
	viper.SetDefault("api.title", "Whodis API")
	viper.SetDefault("api.description", "API to authenticate/authorize discord users")
	viper.SetDefault("api.version", "0.1")
	viper.SetDefault("api.path", "/api/v1")
	viper.SetDefault("api.hostname", "whodis.netsoc.cloud")
	viper.SetDefault("api.secret", "")

	viper.SetDefault("oauth.google.id", "")
	viper.SetDefault("oauth.google.secret", "")

	viper.SetDefault("discord.token", "")
}
