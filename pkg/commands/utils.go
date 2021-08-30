package commands

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func checkDirectMessage(i *discordgo.Interaction) (*discordgo.User, error) {
	if i.GuildID == "" {
		return nil, errors.New("This command is only available from inside a valid server.")
	}
	return i.Member.User, nil
}

func interactionResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, errorMessage string, tagError bool) {
	if tagError {
		errorMessage = fmt.Sprintf("Encountered error: %v", errorMessage)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errorMessage,
			Flags:   1 << 6,
		},
	})
}
