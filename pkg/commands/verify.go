package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/uccnetsoc/whodis/pkg/verify"
)

// VerifyCommand inits the verification process.
func VerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		switch err.(type) {
		case *discordgo.RESTError:
			return "You must run this command in a server.", err
		default:
			return "", err
		}
	}
	user := i.Member.User
	if user == nil {
		user = i.User
	}
	if err = verify.Transition(&verify.StateParams{User: user, Guild: guild}); err != nil {
		return "", err
	}
	return "We have sent you a DM with instruction on how to continue the verification process", nil
}
