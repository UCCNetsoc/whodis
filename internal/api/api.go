package api

import (
	"context"
	"net/http"

	"github.com/Strum355/log"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title Whodis API
// @version 0.1
// @description API to authorize users with given mail domain access to discord guilds
func InitAPI(s *discordgo.Session) {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		ctx := context.WithValue(context.Background(), log.Key, log.Fields{
			"method":  param.Method,
			"path":    param.Path,
			"status":  param.StatusCode,
			"latency": param.Latency,
			"agent":   param.Request.UserAgent(),
		})
		log.WithContext(ctx).Info("invoked request")
		return ""
	}))
	r.Use(gin.Recovery())

	r.GET("/google/login", googleLoginHandler)
	r.GET("/google/auth", googleAuthHandler)

	r.GET("/discord/auth", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/google/login?state="+c.Query("state"))
	})

	r.GET("/verify", createVerifyHandler(s))

	r.GET("/invite", func(c *gin.Context) {
		c.Writer.Header().Add("Location", viper.GetString("discord.bot.invite"))
		c.Writer.WriteHeader(308)
	})

	r.NoRoute(func(c *gin.Context) { infoTemplate.Execute(c.Writer, nil) })

	r.Run(":" + viper.GetString("api.port"))
}
