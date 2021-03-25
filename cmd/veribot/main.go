package main

import (
	"github.com/uccnetsoc/veribot/config"
	"github.com/uccnetsoc/veribot/internal/api"
	"github.com/uccnetsoc/veribot/pkg/models"
)

func main() {
	config.InitConfig()

	api.InitDiscordOAuth()
	api.InitGoogleOAuth()

	models.InitModels()

	api.InitAPI()
}
