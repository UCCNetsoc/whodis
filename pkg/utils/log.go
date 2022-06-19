package utils

import (
	"github.com/Strum355/log"
	"github.com/bwmarrin/discordgo"
)

// SendMessage sends a message to the provided logging channel
func SendLogMessage(s *discordgo.Session, channelId string, message interface{}) {
	_, err := s.ChannelMessageSend(channelId, message.(string))
	if err != nil {
		log.WithError(err)
	}
}
