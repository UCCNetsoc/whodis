package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// VerifyCommand inits the verification process.
func VerifyCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	user := i.Member.User
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		return &interactionError{err, "Failed to get guild from guildID"}
	}
	guildRoleNames := map[string]string{}
	guildRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		return &interactionError{err, "Failed to get guildRoles from guildID"}
	}
	for _, guildRole := range guildRoles {
		guildRoleNames[guildRole.ID] = guildRole.Name
	}
	for _, roleID := range i.Member.Roles {
		roleName, ok := guildRoleNames[roleID]
		if ok && roleName == viper.GetString("discord.member.role") {
			return &interactionError{errors.New("member role is already assigned to user"), "You are already assigned the `" + viper.GetString("discord.member.role") + "` role."}
		}
	}
	uid, err := utils.Encrypt(user.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		return &interactionError{err, "Failed to encrypt userID"}
	}
	gid, err := utils.Encrypt(guild.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		return &interactionError{err, "Failed to encrypt guildID"}
	}
	encoded, err := utils.Encrypt(fmt.Sprintf("%s.%s", uid, gid), []byte(viper.GetString("api.secret")))
	if err != nil {
		return &interactionError{err, "Failed to encrypt user info digest"}
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hey **%s**! Welcome to **%s**!", user.Username, guild.Name),
			Flags:   1 << 6, // Whisper Flag
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Click here to register",
							Style:    discordgo.LinkButton,
							Disabled: false,
							URL:      viper.GetString("api.url") + "/discord/auth?state=" + encoded,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return &interactionError{err, "Unable to respond to interaction"}
	}
	time.AfterFunc(time.Second*15, func() {
		s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Button has timed out, use /verify or click registration button to try again!",
							Style:    discordgo.LinkButton,
							Disabled: true,
							URL:      "https://netsoc.co/rk",
						},
					},
				},
			},
		})
	})
	return nil
}
