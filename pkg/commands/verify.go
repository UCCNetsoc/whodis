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

// VerifyCommand creates a button component with encrypted & signed state specific to the user to go through OAuth flow.
func VerifyCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	user := i.Member.User
	guildID := i.GuildID
	guild, err := s.Guild(guildID)
	if err != nil {
		return &interactionError{err, "Failed to get guild from guildID"}
	}

	guildRoleNames := map[string]string{}
	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		return &interactionError{err, "Failed to get guildRoles from guildID"}
	}

	for _, guildRole := range guildRoles {
		guildRoleNames[guildRole.ID] = guildRole.Name
	}

	for _, roleID := range i.Member.Roles {
		roleName, ok := guildRoleNames[roleID]
		if ok && roleName == viper.GetString("discord.member.role") {
			return &interactionError{
				errors.New("member role is already assigned to user"), "You are already assigned the `" +
					viper.GetString("discord.member.role") + "` role.",
			}
		}
	}

	// encode the userID, guildID, ( welcome channel, logging channel, and roles to give to verified user )
	encoded, err := utils.Encrypt(
		fmt.Sprintf("%s.%s.%s", user.ID, guild.ID, i.MessageComponentData().CustomID[2:]), []byte(viper.GetString("api.secret")),
	)
	if err != nil {
		return &interactionError{err, "Failed to encrypt user info digest"}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Hey **%s**! Welcome to **%s**!\nClick the link below and make sure to sign in with your %s account.",
				user.Username, guild.Name, viper.GetString("oauth.google.domain"),
			),
			Flags: 1 << 6, // Whisper Flag
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: fmt.Sprintf("Click here to register your %s %s", viper.GetString(
								"oauth.google.domain"), " account.",
							),
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
							Label:    "Button timed out. Use /verify or click the registration button again to retry!",
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
