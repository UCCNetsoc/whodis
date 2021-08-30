package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Strum355/log"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/utils"
)

// @title Whodis API
// @version 0.1
// @description API to authorize users with given mail domain access to discord guilds
func InitAPI(s *discordgo.Session) {
	r := gin.Default()

	r.GET("/google/login", googleLoginHandler)
	r.GET("/google/auth", googleAuthHandler)

	r.GET("/discord/auth", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/google/login?state=%s", c.Query("state")))
	})

	r.GET("/verify", func(c *gin.Context) {
		encodedDigest := c.Query("state")
		if len(encodedDigest) == 0 {
			log.Error("Error parsing discord digest")
			c.Writer.WriteString("Error parsing discord digest")
			c.Status(http.StatusInternalServerError)
			return
		}
		decodedDigest, err := utils.Decrypt(encodedDigest, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.WithError(err).Error("Error decoding discord digest")
			c.Writer.WriteString("Error decoding discord digest")
			c.Status(http.StatusInternalServerError)
			return
		}
		encodedData := strings.Split(decodedDigest, ".")
		encodedUID, encodedGID := encodedData[0], encodedData[1]
		decodedUID, err := utils.Decrypt(encodedUID, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.WithError(err).Error("Error decoding discord userID")
			c.Writer.WriteString("Error decoding discord userID")
			c.Status(http.StatusInternalServerError)
			return
		}
		decodedGID, err := utils.Decrypt(encodedGID, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.WithError(err).Error("Error decoding discord guildID")
			c.Writer.WriteString("Error decoding discord guildID")
			c.Status(http.StatusInternalServerError)
			return
		}
		roleID := ""
		roles, _ := s.GuildRoles(decodedGID)
		for _, role := range roles {
			if role.Name == "Member" {
				roleID = role.ID
			}
		}
		if roleID == "" {
			log.Error("Error finding `Member` role")
			c.Writer.WriteString("Error finding `Member` role")
			c.Status(http.StatusInternalServerError)
			return
		}
		err = s.GuildMemberRoleAdd(decodedGID, decodedUID, roleID)
		if err != nil {
			log.WithError(err).Error("Error adding `Member` role to user")
			c.Writer.WriteString("Error adding `Member` role to user")
			c.Status(http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("\n\t`Member` role has been added. user='%s' guild='%s'", decodedUID, decodedGID))
		c.Status(200)
		c.Writer.WriteString("Success")
	})
	r.Run(":8080")
}
