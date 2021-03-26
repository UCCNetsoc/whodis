package commands

import (
	"github.com/bwmarrin/discordgo"
)

// VerifyCommand inits the verification process.
func VerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "We have sent you a DM with instruction on how to continue the verification process",
		},
	})
}
