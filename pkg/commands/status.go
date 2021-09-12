package commands

import (
	"context"
	"errors"

	"github.com/Strum355/log"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const (
	statusOK    string = "✅"
	statusError string = "❌"
)

// StatusCommand checks for server compatibility issues.
func StatusCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError {
	if !memberPermissionCheck(s, i.GuildID, i.Member) {
		return &interactionError{errors.New("user has invalid permissions"), "You do not have valid permissions to use this command"}
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
	roleExists, roleAccess := botPermissionCheck(s, i.GuildID)
	return &discordgo.MessageEmbed{
		Title: "Whodis Setup Status",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Member role exists in guild:",
				Value: statusEmote[roleExists],
			},
			{
				Name:  "Permission to manage Member role:",
				Value: statusEmote[roleAccess],
			},
		},
	}
}

func generateRoleMap(s *discordgo.Session, guildID string) map[string]*discordgo.Role {
	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		log.WithError(err).Error("Failed to generate role map")
		return nil
	}
	guildRoleMap := map[string]*discordgo.Role{}
	for _, role := range guildRoles {
		guildRoleMap[role.ID] = role
	}
	return guildRoleMap
}

func memberPermissionCheck(s *discordgo.Session, guildID string, guildMember *discordgo.Member) bool {
	guildRoleMap := generateRoleMap(s, guildID)
	if guildRoleMap == nil {
		return false
	}
	guild, err := s.Guild(guildID)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return false
	}
	if guildMember.User.ID == guild.OwnerID {
		return true // User is guild owner
	}
	for _, guildMemberRoleID := range guildMember.Roles {
		memberRole := guildRoleMap[guildMemberRoleID]

		memberRoleManageRole := memberRole.Permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles
		memberRoleAdmin := memberRole.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator

		if memberRoleManageRole || memberRoleAdmin {
			return true
		}
	}
	return false
}

func botPermissionCheck(s *discordgo.Session, guildID string) (roleExists bool, roleAccess bool) {
	guildRoleMap := generateRoleMap(s, guildID)
	if guildRoleMap == nil {
		return false, false
	}
	var memberRole *discordgo.Role
	for _, role := range guildRoleMap {
		if role.Name == viper.GetString("discord.member.role") {
			memberRole = role
		}
	}
	if memberRole == nil {
		return false, false
	}
	bot, err := s.GuildMember(guildID, viper.GetString("discord.app.id"))
	if err != nil {
		log.WithError(err).Error(err.Error())
		return true, false
	}
	for _, botRoleID := range bot.Roles {
		botRole := guildRoleMap[botRoleID]

		botRoleManageRole := botRole.Permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles
		botRoleAdmin := botRole.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator

		botRoleHigherRank := botRole.Position > memberRole.Position

		if (botRoleManageRole || botRoleAdmin) && botRoleHigherRank {
			return true, true
		}
	}
	return true, false
}
