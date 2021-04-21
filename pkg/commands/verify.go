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
		switch err.(type) {
		case *discordgo.RESTError:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: "You must run this command in a server.",
				},
			})
			log.Printf("User invalid guild %T", err)
			return
		default:
			log.Println(err)
			return
		}
	}
	if err = verify.Transition(&verify.StateParams{User: i.User, Guild: guild}); err != nil {
		log.Println(err)
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "We have sent you a DM with instruction on how to continue the verification process",
		},
	})
}
