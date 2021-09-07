package commands

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
)

const (
	statusOK    string = "✅"
	statusError string = "❌"
)

// StatusCommand checks for server compatibility issues.
func StatusCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	if !permissionCheck(s, i.GuildID, i.Member.User.ID) {
		return &interactionError{errors.New("User has invalid permissions"), "You do not have valid permissions to use this command"}
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{createStatusEmbed(s, i)},
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return &interactionError{err, err.Error()}
	}
	return nil
}

func createStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	statusEmote := map[bool]string{true: statusOK, false: statusError}
	return &discordgo.MessageEmbed{
		Title: "Whodis Setup Status",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Permission to manage roles:",
				Value: statusEmote[permissionCheck(s, i.GuildID, i.Member.User.ID)],
			},
			{
				Name:  "Access to Member role:",
				Value: statusEmote[roleCheck(s, i)],
			},
		},
	}
}

func permissionCheck(s *discordgo.Session, guildID string, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}
	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		return false
	}
	guildRolePerms := map[string]int64{}

	for _, guildRole := range guildRoles {
		guildRolePerms[guildRole.ID] = guildRole.Permissions
	}
	for _, memberRole := range member.Roles {
		rolePerms, ok := guildRolePerms[memberRole]
		if ok && (rolePerms&discordgo.PermissionManageRoles != 0) {
			return true
		}
	}
	return member.Permissions&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles
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
