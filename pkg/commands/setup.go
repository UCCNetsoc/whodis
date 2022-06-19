package commands

import (
	"context"
	"errors"

	"github.com/UCCNetsoc/whodis/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

// SetupCommand creates a component in the specified channel to init registration.
func SetupCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		return &interactionError{message: "User must be admin to run /setup.", err: errors.New("user must be admin to run /setup")}
	}
	g, err := s.Guild(i.GuildID)
	if err != nil {
		return &interactionError{err: err, message: "Unable to query guild"}
	}
	m := &discordgo.MessageSend{
		Content: "Welcome to **" + g.Name + "**!",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Click here to register.",
						Style:    discordgo.PrimaryButton,
						CustomID: "v" + concatData(ctx, s, i),
					},
				},
			},
		},
	}
	if _, err := s.ChannelMessageSendComplex(i.ChannelID, m); err != nil {
		return &interactionError{err: err, message: "Unable to send message to allow registration."}
	}
	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6, // Whisper Flag
			Content: "This channel has been setup as the welcome room.\nMake sure users don't have send message permissions on this channel to avoid register message being buried by spam.",
		},
	}

	utils.SendLogMessage(s, i.ApplicationCommandData().Options[1].ChannelValue(s).ID,
		"This channel has been setup as the logging channel.\nDue to the stateless and privacy-conscious nature of Whodis, some errors may not get logged. It is advisable to keep a help channel for new members to who can't authenticate, in case logs don't reflect the errors.",
	)

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return &interactionError{err: err, message: "Couldn't reply to interaction."}
	}
	return nil
}

func concatData(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) string {
	data := "." + i.ApplicationCommandData().Options[0].ChannelValue(s).ID
	data += "." + i.ApplicationCommandData().Options[1].ChannelValue(s).ID
	for _, role := range i.ApplicationCommandData().Options[2:] {
		data += "." + role.RoleValue(s, i.GuildID).ID
	}
	return data
}
