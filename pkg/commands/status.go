package commands

import (
	"context"

	"github.com/Strum355/log"

	"github.com/bwmarrin/discordgo"
)

const (
	statusOK    string = "✅"
	statusError string = "❌"
)

// StatusCommand checks for server compatibility issues.
func StatusCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	p, err := s.State.UserChannelPermissions(user.ID, i.ChannelID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Could not get user permissions.")
		interactionResponseError(s, i, err.Error(), true)
		return
	}
	if p&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		log.Error("User has invalid permissions")
		interactionResponseError(s, i, "You do not have valid permissions to use this command.", true)
		return
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{createStatusEmbed(s, i)},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.WithError(err)
		interactionResponseError(s, i, err.Error(), true)
		return
	}
}

func createStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	statusEmote := map[bool]string{true: statusOK, false: statusError}
	return &discordgo.MessageEmbed{
		Title: "Whodis Setup Status",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Permission to manage roles:",
				Value: statusEmote[permissionCheck(s, i)],
			},
			{
				Name:  "Access to Member role:",
				Value: statusEmote[roleCheck(s, i)],
			},
		},
	}
}

func permissionCheck(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	p, err := s.State.UserChannelPermissions(s.State.User.ID, i.ChannelID)
	if err != nil {
		return false
	}
	if p&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		return false
	}
	return true
}

func roleCheck(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	roles, _ := s.GuildRoles(i.GuildID)
	for _, role := range roles {
		if role.Name == "Member" {
			return true
		}
	}
	return false
}
