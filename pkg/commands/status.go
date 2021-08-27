package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// StatusCommand checks for server compatibility issues.
func StatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	// guild, _ := s.Guild(i.GuildID)
	if user == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command is only available from inside a valid server.",
				Flags:   1 << 6,
			},
		})
		return
	}
	p, err := s.State.UserChannelPermissions(user.ID, i.ChannelID)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error encountered: " + err.Error(),
				Flags:   1 << 6,
			},
		})
		return
	}
	if p&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		log.Println("User has invalid permissions")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error encountered: You do not have valid permissions to use this command.",
				Flags:   1 << 6,
			},
		})
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
		log.Println(err)
		return
	}
}

func createStatusEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) (emb *discordgo.MessageEmbed) {
	permissionCheck, roleCheck := "✅", "❌"

	p, err := s.State.UserChannelPermissions(s.State.User.ID, i.ChannelID)
	if err != nil {
		permissionCheck = "❌"
	}
	if p&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		permissionCheck = "❌"
	}
	roles, _ := s.GuildRoles(i.GuildID)
	for _, role := range roles {
		if role.Name == "Member" {
			roleCheck = "✅"
		}
	}
	emb = &discordgo.MessageEmbed{
		Title: "Server Status",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Permission to manage roles:",
				Value: permissionCheck,
				// Inline: true,
			},
			{
				Name:  "Access to Member role:",
				Value: roleCheck,
				// Inline: true,
			},
		},
	}
	return
}
