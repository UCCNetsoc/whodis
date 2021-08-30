package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Strum355/log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/config"
	"github.com/uccnetsoc/whodis/internal/api"
	"github.com/uccnetsoc/whodis/pkg/commands"
)

func main() {
	log.InitSimpleLogger(&log.Config{Output: os.Stdout})
	config.InitConfig()

	api.InitGoogleOAuth()

	s, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.WithError(err)
		return
	}
	s.StateEnabled = true
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	go api.InitAPI(s)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info("Bot is starting")
	})
	s.Open()
	commands.RegisterSlashCommands(s)
	log.Info("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Cleanly exiting")
	s.Close()
}
