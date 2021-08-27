package commands

import (
	"github.com/bwmarrin/discordgo"
)

// RegisterSlashCommands adds all slash commands to the session.
func RegisterSlashCommands(s *discordgo.Session) {
	commands.Add(
		&discordgo.ApplicationCommand{
			Name:        "verify",
			Description: "Start the verification process to get into the current server.",
		},
		VerifyCommand,
	)
	commands.Add(
		&discordgo.ApplicationCommand{
			Name:        "status",
			Description: "Verify that the Whodis bot and your server are both in working order.",
		},
		StatusCommand,
	)
	commands.Register(s)
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

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
		if handler, ok := commands.handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
	for _, comm := range c.commands {
		if _, err := s.ApplicationCommandCreate("879010291126517810", "875053603012870215", comm); err != nil {
			return err
		}
	}
	return nil
}
