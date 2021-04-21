package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error)

type Commands struct {
	commands []*discordgo.ApplicationCommand
	handlers map[string]CommandHandler
}

var (
	commands = &Commands{}
)

// Add a commands to the slash commands.
func (c *Commands) Add(com *discordgo.ApplicationCommand, handler CommandHandler) {
	c.commands = append(c.commands, com)
	if c.handlers == nil {
		c.handlers = map[string]CommandHandler{}
	}
	c.handlers[com.Name] = handler
}

// Register all slash commands.
func (c *Commands) Register(s *discordgo.Session) error {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commands.handlers[i.Data.Name]; ok {
			resp, err := handler(s, i)
			if err != nil {
				log.Println(err)
			}
			if resp != "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: resp,
					},
				})
			}
		}
	})
	for _, comm := range c.commands {
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, "", comm); err != nil {
			return err
		}
	}
	return nil
}

// RegisterSlashCommands adds all slash commands to the session.
func RegisterSlashCommands(s *discordgo.Session) {
	commands.Add(
		&discordgo.ApplicationCommand{
			Name:        "verify",
			Description: "Start the verification process to get into the current server.",
		},
		VerifyCommand,
	)
	commands.Register(s)
}
