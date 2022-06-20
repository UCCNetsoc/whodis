package api

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Strum355/log"
	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const success = "You have successfully registered for the discord server, you can return to discord, where you have elevated permissions."

//go:embed assets/*
var pageData embed.FS
var infoTemplate = template.Must(template.ParseFS(pageData, "assets/info.html"))
var resultTemplate = template.Must(template.ParseFS(pageData, "assets/result.html"))

// createVerifyHandler creates handles the /verify callback endpoint.
func createVerifyHandler(s *discordgo.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, gid, announce_cid, logging_cid, rid, err := dataFromState(c.Query("state"))
		if err != nil {
			resultTemplate.Execute(c.Writer,
				AccessErrorResponse(http.StatusInternalServerError, "Error decoding state from URL.", err),
			)
			log.WithError(err).Error("Error decoding state from URL.")
			utils.SendLogMessage(s, logging_cid, "Failed to verify integrity of unique link for user. If issues persist, create an issue at https://github.com/UCCNetsoc/whodis")
			return
		}

		user, err := s.User(uid)
		if err != nil {
			resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusInternalServerError, "Error getting user", err))
			utils.SendLogMessage(s, logging_cid, "Failed to get user `"+uid+"` from Discord API to add roles.")
			return
		}

		if err := addRoles(s, gid, uid, rid); err != nil {
			resultTemplate.Execute(c.Writer,
				AccessErrorResponse(http.StatusInternalServerError, "Error adding roles to Discord user.", err),
			)

			log.Error("Error adding roles to Discord user: " + user.Username + "#" + user.Discriminator + " on guild: " + gid)

			utils.SendLogMessage(s, logging_cid,
				fmt.Sprintf("Failed to add roles to user %s. Ensure that Whodis has permission to manage user roles.", user.Username+"#"+user.Discriminator),
			)
			return
		}

		// Get welcome channel.
		announceChannelID, ok := viper.GetStringMapString("discord.guild.members.channel")[gid]
		if !ok {
			if announce_cid != "" {
				announceChannelID = announce_cid
			} else if announceChannelID, err = getDefaultChannel(s, gid); err != nil {
				resultTemplate.Execute(c.Writer,
					AccessErrorResponse(http.StatusInternalServerError, "Error querying channels for welcome message.", err),
				)
				utils.SendLogMessage(s, logging_cid,
					"Failed to get announce channel, ensure the channel is a text channel and try using the **/setup** command again.",
				)
				return
			}
		}
		if announceChannelID == "" {
			resultTemplate.Execute(c.Writer, *AccessSuccessResponse(success, uid, gid))
			return
		}

		// Get welcome message.
		welcomeMessage, ok := viper.GetStringMapString("discord.guild.members.welcome")[gid]
		if !ok {
			welcomeMessage = viper.GetString("discord.welcome.default")
		}

		// Welcome user in channel.
		if _, err := s.ChannelMessageSend(announceChannelID, "Welcome "+user.Mention()+"! "+welcomeMessage); err != nil {
			resultTemplate.Execute(c.Writer,
				AccessErrorResponse(http.StatusInternalServerError, "Error sending message to welcome channel", err),
			)
			utils.SendLogMessage(s, logging_cid, "Failed to send message to announce channel, ensure permissions allow the bot to do so, or try the **/setup** command again.")
			return
		}
		resultTemplate.Execute(c.Writer, *AccessSuccessResponse(success, uid, gid))
		log.Info("Successfully verified user: " + user.Username + "#" + user.Discriminator + " on guild: " + gid)
		roles := []string{"**" + viper.GetString("discord.member.role") + "**"}
		for _, role_id := range rid {
			role, err := s.State.Role(gid, role_id)
			if err != nil {
				utils.SendLogMessage(s, logging_cid, "Failed to get role `"+role_id+"` from Discord API to add to user.")
				continue
			}
			roles = append(roles, "**"+role.Name+"**")
		}
		utils.SendLogMessage(s, logging_cid,
			fmt.Sprintf("User %s successfully verified, roles added: %s", user.Mention(), strings.Join(roles, ", ")),
		)
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

func addRoles(s *discordgo.Session, gid, uid string, rid []string) error {
	roles, err := s.GuildRoles(gid)
	if err != nil {
		return err
	}
	roleID := utils.GetRoleIDFromName(roles, viper.GetString("discord.member.role"))
	if roleID == "" {
		return errors.New("no role called: " + viper.GetString("discord.member.role"))
	}
	if err := s.GuildMemberRoleAdd(gid, uid, roleID); err != nil {
		return err
	}
	for _, additionalRoleID := range rid {
		if err := s.GuildMemberRoleAdd(gid, uid, additionalRoleID); err != nil {
			return err
		}
	}
	return nil
}

func dataFromState(state string) (uid string, gid string, announce_cid string, logging_cid string, rid []string, err error) {
	if len(state) == 0 {
		err = errors.New("no state found")
		return
	}
	var decodedDigest string
	if decodedDigest, err = utils.Decrypt(state, []byte(viper.GetString("api.secret"))); err != nil {
		return
	}
	log.Info("Decoded digest: " + decodedDigest)
	encodedData := strings.Split(decodedDigest, ".")
	uid, gid = encodedData[0], encodedData[1]
	announce_cid = encodedData[2]
	logging_cid = encodedData[3]
	if len(encodedData) > 3 {
		rid = encodedData[4:]
	}
	return
}
