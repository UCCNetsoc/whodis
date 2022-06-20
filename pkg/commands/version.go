package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func VersionCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	version := viper.GetString("bot.version")
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Whodis Version",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Version",
							Value:  fmt.Sprintf("[%s](https://github.com/UCCNetsoc/whodis/releases/tag/%s)", version, version),
							Inline: true,
						},
					},
				},
			},
		},
	}
	if err := s.InteractionRespond(i.Interaction, response); err != nil {
		return &interactionError{err, err.Error()}
	}

	return nil
}
