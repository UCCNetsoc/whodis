package api

import (
	"embed"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const success = "You have successfully registered for the discord server."

//go:embed assets/*
var pageData embed.FS
var infoTemplate = template.Must(template.ParseFS(pageData, "assets/info.html"))
var resultTemplate = template.Must(template.ParseFS(pageData, "assets/result.html"))

// createVerifyHandler creates handles the /verify callback endpoint.
func createVerifyHandler(s *discordgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, gid, err := idsFromState(c.Query("state"))
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error decoding state from URL.", err))
			return
		}
		if err := addRoles(s, gid, uid); err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error adding roles to Discord user.", err))
			return
		}

		// Get welcome channel.
		channelID, ok := viper.GetStringMapString("discord.guild.members.channel")[gid]
		if !ok {
			if channelID, err = getDefaultChannel(s, gid); err != nil {
				resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error querying channels for welcome message.", err))
				return
			}
		}
		if channelID == "" {
			resultTemplate.Execute(c.Writer, *AccessSuccessResponse(success, uid, gid))
			return
		}

		// Welcome user in channel.
		user, err := s.User(uid)
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error getting user", err))
			return
		}
		if _, err := s.ChannelMessageSend(channelID, "Welcome **"+user.Mention()+"**! Thanks for registering!"); err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error sending message to welcome channel", err))
			return
		}
		resultTemplate.Execute(c.Writer, *AccessSuccessResponse(success, uid, gid))
	}
}

func getDefaultChannel(s *discordgo.Session, gid string) (string, error) {
	var channelID string
	channels, err := s.GuildChannels(gid)
	if err != nil {
		return "", err
	}
	for _, channel := range channels {
		if channel.Name == viper.GetString("discord.channel.default") {
			channelID = channel.ID
			break
		}
	}
	return channelID, nil
}

func addRoles(s *discordgo.Session, gid, uid string) error {
	roles, err := s.GuildRoles(gid)
	if err != nil {
		return err
	}
	roleID := utils.GetRoleIDFromName(roles, viper.GetString("discord.member.role"))
	if roleID == "" {
		return errors.New("no role called: " + viper.GetString("discord.member.role"))
	}
	for _, roleName := range viper.GetStringSlice("discord.additional.roles") {
		id := utils.GetRoleIDFromName(roles, roleName)
		if id == "" {
			continue
		}
		if err := s.GuildMemberRoleAdd(gid, uid, id); err != nil {
			return err
		}
	}
	if err := s.GuildMemberRoleAdd(gid, uid, roleID); err != nil {
		return err
	}
	return nil
}

func idsFromState(state string) (uid string, gid string, err error) {
	if len(state) == 0 {
		err = errors.New("no state found")
		return
	}
	var decodedDigest string
	if decodedDigest, err = utils.Decrypt(state, []byte(viper.GetString("api.secret"))); err != nil {
		return
	}
	encodedData := strings.Split(decodedDigest, ".")
	encodedUID, encodedGID := encodedData[0], encodedData[1]
	if uid, err = utils.Decrypt(encodedUID, []byte(viper.GetString("api.secret"))); err != nil {
		return
	}
	if gid, err = utils.Decrypt(encodedGID, []byte(viper.GetString("api.secret"))); err != nil {
		return
	}
	return
}
