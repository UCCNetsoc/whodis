package config

import "github.com/spf13/viper"

func initDefaults() {
	viper.SetDefault("discord.token", "")
	viper.SetDefault("discord.app.id", "")
	viper.SetDefault("discord.member.role", "Member")
	viper.SetDefault("discord.additional.role", "")
	viper.SetDefault("discord.bot.invite", "")
	viper.SetDefault("discord.guild.members.channel", map[string]string{})
	viper.SetDefault("discord.channel.default", "general")
	viper.SetDefault("discord.guild.members.welcome", map[string]string{})
	viper.SetDefault("discord.welcome.default", "Thanks for registering!")

	viper.SetDefault("api.url", "")
	viper.SetDefault("api.port", "")
	viper.SetDefault("api.secret", "")

	viper.SetDefault("oauth.google.id", "")
	viper.SetDefault("oauth.google.secret", "")
	viper.SetDefault("oauth.google.domain", "umail.ucc.ie")

}
