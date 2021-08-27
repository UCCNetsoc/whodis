package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/utils"
)

// @title Veribot API
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
			log.Println("Error parsing discord digest")
			c.Writer.WriteString("Error parsing discord digest")
			c.Status(http.StatusInternalServerError)
			return
		}
		decodedDigest, err := utils.Decrypt(encodedDigest, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.Println("Error decoding discord digest\n", err)
			c.Writer.WriteString("Error decoding discord digest")
			c.Status(http.StatusInternalServerError)
			return
		}
		encodedData := strings.Split(decodedDigest, ".")
		encodedUID, encodedGID := encodedData[0], encodedData[1]
		decodedUID, err := utils.Decrypt(encodedUID, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.Println("Error decoding discord userID\n", err)
			c.Writer.WriteString("Error decoding discord userID")
			c.Status(http.StatusInternalServerError)
			return
		}
		decodedGID, err := utils.Decrypt(encodedGID, []byte(viper.GetString("api.secret")))
		if err != nil {
			log.Println("Error decoding discord guildID\n", err)
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
			log.Println("Error finding `Member` role")
			c.Writer.WriteString("Error finding `Member` role")
			c.Status(http.StatusInternalServerError)
			return
		}
		s.GuildMemberRoleAdd(decodedGID, decodedUID, roleID)
		c.Status(200)
		c.Writer.WriteString("Success")
	})
	r.Run(":8080")
}
