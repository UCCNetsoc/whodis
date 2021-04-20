package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/uccnetsoc/whodis/pkg/verify"
)

// VerifyCommand inits the verification process.
func VerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println(err)
	}
	if err = verify.Transition(&verify.StateParams{User: i.User, Guild: guild}); err != nil {
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "We have sent you a DM with instruction on how to continue the verification process",
		},
	})
}
