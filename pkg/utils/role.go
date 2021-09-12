package utils

import "github.com/bwmarrin/discordgo"

// GetRoleIDFromName finds a discord role by name and returns ID.
func GetRoleIDFromName(roles []*discordgo.Role, name string) string {
	var roleID string
	for _, role := range roles {
		if role.Name == name {
			roleID = role.ID
			break
		}
	}
	return roleID
}
