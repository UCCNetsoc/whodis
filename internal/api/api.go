package api

import (
	"net/http"
	"strings"

	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

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
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error parsing discord digest", nil))
			return
		}
		decodedDigest, err := utils.Decrypt(encodedDigest, []byte(viper.GetString("api.secret")))
		if err != nil {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord digest", err))
			return
		}
		encodedData := strings.Split(decodedDigest, ".")
		encodedUID, encodedGID := encodedData[0], encodedData[1]
		decodedUID, err := utils.Decrypt(encodedUID, []byte(viper.GetString("api.secret")))
		if err != nil {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord userID", err))
			return
		}
		decodedGID, err := utils.Decrypt(encodedGID, []byte(viper.GetString("api.secret")))
		if err != nil {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error decoding discord guildID", err))
			return
		}
		roleID := ""
		roles, err := s.GuildRoles(decodedGID)
		if err != nil {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error getting guild roles", err))
			return
		}
		for _, role := range roles {
			if role.Name == "Member" {
				roleID = role.ID
				break
			}
		}
		if roleID == "" {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error finding Member role", nil))
			return
		}
		err = s.GuildMemberRoleAdd(decodedGID, decodedUID, roleID)
		if err != nil {
			c.JSON(AccessErrorResponse(http.StatusInternalServerError, "Error adding Member role to user", err))
			return
		}
		c.JSON(AccessSuccessResponse("Role has been added to user", decodedUID, decodedGID, roleID))
	})
	r.Run(":" + viper.GetString("api.port"))
}
