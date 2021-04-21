package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/config"
	"github.com/uccnetsoc/whodis/internal/api"
	"github.com/uccnetsoc/whodis/pkg/commands"
	"github.com/uccnetsoc/whodis/pkg/models"
	"github.com/uccnetsoc/whodis/pkg/verify"
)

func main() {
	config.InitConfig()

	api.InitDiscordOAuth()
	api.InitGoogleOAuth()

	models.InitModels()

	s, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.Fatal(err)
		return
	}
	go api.InitAPI(s)
	verify.Init(s)
	s.Open()
	commands.RegisterSlashCommands(s)
	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Cleanly exiting")
	s.Close()

}
