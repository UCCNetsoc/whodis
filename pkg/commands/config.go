package commands

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/bwmarrin/discordgo"
	"github.com/uccnetsoc/whodis/pkg/models"
)

var configOptions = []*discordgo.ApplicationCommandOption{
	{
		Name:        "roles",
		Description: "List of roles to be added when user is authenticated.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	},
	{
		Name:        "domains",
		Description: "List of email domains google accounts being verified can have. Leave blank for all domains. Can be separeted by a space and/or comma.",
		Type:        discordgo.ApplicationCommandOptionSubCommand,
	},
}

func ConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
	for _, sub := range i.Data.Options {
		switch sub.Name {
		case "set":
			if len(sub.Options) == 0 {
				return "Must provide value to set.", nil
			}
			var (
				value interface{}
				ok    bool
			)
			if value, ok = sub.Options[0].Value.(string); !ok {
				return "", fmt.Errorf("invalid data type %T for config set value", value)
			}
			name := sub.Options[0].Name
			switch name {
			case "domains":
				var userErr string
				if value, userErr = parseDomains(value.(string)); userErr != "" {
					return userErr, nil
				}
			case "roles":
				var userErr string
				if value, userErr = parseList(value.(string)); userErr != "" {
					return userErr, nil
				}
			}
			if err := models.DBClient.SetConfigItem(i.GuildID, name, value); err != nil {
				return "Internal error occured.", err
			}
			return fmt.Sprintf("Config parameter `%s` set.", name), nil
		}
	}
	return "No parameter provided.", nil
}

func parseDomains(value string) ([]string, string) {
	return parseList(value, func(value string) string {
		if !govalidator.IsDNSName(value) {
			return fmt.Sprintf("%s is not a valid DNS name", value)
		}
		return ""
	})
}

func parseList(value string, filter ...func(value string) string) ([]string, string) {
	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, ",", " ")
	}
	fields := strings.Fields(value)
	for _, field := range fields {
		if len(filter) > 0 {
			if strErr := filter[0](field); strErr != "" {
				return nil, strErr
			}
		}
	}
	return fields, ""
}
