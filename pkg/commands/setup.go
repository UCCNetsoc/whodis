package commands

import (
	"context"
	"errors"

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
		Content: "Welcome @everyone to **" + g.Name + "**!",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Click here to register.",
						Style:    discordgo.PrimaryButton,
						CustomID: "verify",
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
			Content: "This channel is now setup as the welcome room.\nMake sure users don't have send message permissions on this channel to avoid register message being buried by spam.",
		},
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return &interactionError{err: err, message: "Couldn't reply to interaciton."}
	}
	return nil
}
