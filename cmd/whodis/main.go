package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Strum355/log"

	"github.com/UCCNetsoc/whodis/config"
	"github.com/UCCNetsoc/whodis/internal/api"
	"github.com/UCCNetsoc/whodis/pkg/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var production *bool

func main() {
	production = flag.Bool("p", false, "enables production with json logging")
	flag.Parse()
	if *production {
		log.InitJSONLogger(&log.Config{Output: os.Stdout})
	} else {
		log.InitSimpleLogger(&log.Config{Output: os.Stdout})
	}

	config.InitConfig()

	api.InitGoogleOAuth()

	s, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.WithError(err)
		return
	}
	s.StateEnabled = true
	go api.InitAPI(s)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info("Bot is has registered handlers")
	})
	s.Open()
	commands.RegisterSlashCommands(s)
	log.Info("Bot is initialising")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info("Cleanly exiting")
	s.Close()
}
