package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/veribot/config"
	"github.com/uccnetsoc/veribot/internal/api"
	"github.com/uccnetsoc/veribot/pkg/commands"
	"github.com/uccnetsoc/veribot/pkg/models"
)

func main() {
	config.InitConfig()

	api.InitDiscordOAuth()
	api.InitGoogleOAuth()

	models.InitModels()

	api.InitAPI()
	s, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Open()
	commands.RegisterSlashCommands(s)
	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Cleanly exiting")
	s.Close()

}
