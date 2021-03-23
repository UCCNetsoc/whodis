package main

import (
	"github.com/uccnetsoc/veribot/config"
	"github.com/uccnetsoc/veribot/internal/api"
)

func main() {
	config.InitConfig()
	api.InitAPI()
}
