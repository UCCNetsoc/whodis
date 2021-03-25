package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uccnetsoc/veribot/docs"
	"github.com/uccnetsoc/veribot/pkg/models"
)

// @title Veribot API
// @version 0.1
// @description API to authorize users with given mail domain access to discord guilds
func InitAPI() {
	docs.SwaggerInfo.Title = viper.GetString("api.title")
	docs.SwaggerInfo.Description = viper.GetString("api.description")
	docs.SwaggerInfo.Version = viper.GetString("api.version")
	docs.SwaggerInfo.BasePath = viper.GetString("api.path")
	docs.SwaggerInfo.Host = viper.GetString("api.hostname")

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/verify", discordLoginHandler)
	r.GET("/discord/auth", discordAuthHandler)

	r.GET("/google/login", middleware(), googleLoginHandler)
	r.GET("/google/auth", googleAuthHandler)

	r.POST("/verify", func(c *gin.Context) {
		var user *models.User
		body, err := ioutil.ReadAll(c.Request.Body)
		log.Println(string(body))
		if err != nil {
			log.Println("ERROR parsing json body\n", err)
			c.String(http.StatusBadRequest, "bad request body")
			return
		}
		json.Unmarshal(body, &user)
		log.Println(user)
		if user.MailDomain != "" && user.MailDomain == viper.GetString("mail.domain") {
			// call discord bot to update role for discord id
			c.String(http.StatusOK, "successfully authorized to access resource")
			log.Println("nice")
			return
		}
	})

	r.Run(":8080")
}
