package api

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//go:embed assets/*
var pageData embed.FS
var infoTemplate = template.Must(template.ParseFS(pageData, "assets/info.html"))
var resultTemplate = template.Must(template.ParseFS(pageData, "assets/result.html"))

// @title Whodis API
// @version 0.1
// @description API to authorize users with given mail domain access to discord guilds
func InitAPI(s *discordgo.Session) {
	r := gin.Default()
	r.GET("/google/login", googleLoginHandler)
	r.GET("/google/auth", googleAuthHandler)

	r.GET("/discord/auth", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/google/login?state="+c.Query("state"))
	})

	r.GET("/verify", func(c *gin.Context) {
		encodedDigest := c.Query("state")
		if len(encodedDigest) == 0 {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error parsing discord digest", nil))
			return
		}
		decodedDigest, err := utils.Decrypt(encodedDigest, []byte(viper.GetString("api.secret")))
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord digest", err))
			return
		}
		encodedData := strings.Split(decodedDigest, ".")
		encodedUID, encodedGID := encodedData[0], encodedData[1]
		decodedUID, err := utils.Decrypt(encodedUID, []byte(viper.GetString("api.secret")))
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord userID", err))
			return
		}
		decodedGID, err := utils.Decrypt(encodedGID, []byte(viper.GetString("api.secret")))
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord guildID", err))
			return
		}
		roleID := ""
		roles, err := s.GuildRoles(decodedGID)
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error getting guild roles", err))
			return
		}
		for _, role := range roles {
			if role.Name == "Member" {
				roleID = role.ID
				break
			}
		}
		if roleID == "" {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error finding Member role", nil))
			return
		}
		err = s.GuildMemberRoleAdd(decodedGID, decodedUID, roleID)
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error adding Member role to user", err))
			return
		}
		resultTemplate.Execute(c.Writer, *AccessSuccessResponse("Role has been added to user", decodedUID, decodedGID, roleID))
	})

	r.GET("/invite", func(c *gin.Context) {
		c.Writer.Header().Add("Location", viper.GetString("discord.bot.invite"))
		c.Writer.WriteHeader(308)
	})

	r.NoRoute(func(c *gin.Context) { infoTemplate.Execute(c.Writer, nil) })

	r.Run(":" + viper.GetString("api.port"))
}
