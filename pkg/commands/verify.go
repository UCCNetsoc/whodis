package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Strum355/log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/utils"
)

// VerifyCommand inits the verification process.
func VerifyCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	guild, _ := s.Guild(i.GuildID)
	for _, roleID := range i.Member.Roles {
		role, _ := s.State.Role(i.GuildID, roleID)
		if role.Name == "Member" {
			log.WithContext(ctx).Error("`Member` role is already assigned to user")
			interactionResponseError(s, i, "This user is already assigned the `Member` role.", false)
			return
		}
	}
	uid, err := utils.Encrypt(user.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to encrypt userID")
		interactionResponseError(s, i, "Failed to encrypt userID", true)
		return
	}
	gid, err := utils.Encrypt(guild.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to encrypt guildID")
		interactionResponseError(s, i, "Failed to encrypt guildID", true)
		return
	}
	encoded, err := utils.Encrypt(fmt.Sprintf("%s.%s", uid, gid), []byte(viper.GetString("api.secret")))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to encrypt user info digest")
		interactionResponseError(s, i, "Failed to encrypt user info digest", true)
		return
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hey **%s**! Welcome to **%s**!", user.Username, guild.Name),
			Flags:   1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Click here to register",
							Style:    discordgo.LinkButton,
							Disabled: false,
							URL:      viper.GetString("api.host") + "/discord/auth?state=" + encoded,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Unable to respond to interaction")
		interactionResponseError(s, i, "Unable to respond to interaction", true)
		return
	}
	time.AfterFunc(time.Second*15, func() {
		s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Button has timed out, use /verify to try again!",
							Style:    discordgo.LinkButton,
							Disabled: true,
							URL:      "https://netsoc.co/rk",
						},
					},
				},
			},
		})
	})
}
