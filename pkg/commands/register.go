package commands

import (
	"context"

	"github.com/Strum355/log"

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

type CommandHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)

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
		callHandler(s, i)
	})
	for _, comm := range c.commands {
		if _, err := s.ApplicationCommandCreate("879010291126517810", "875053603012870215", comm); err != nil {
			return err
		}
	}
	return nil
}

// Call command handler.
func callHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	commandAuthor, err := checkDirectMessage(i.Interaction)
	if err != nil {
		log.WithError(err).Error("Failed to invoke command")
		interactionResponseError(s, i, err.Error(), true)
		return
	}
	pCheck, rCheck := statusCheck(s, i)
	if !pCheck || !rCheck {
		log.Error("Invalid bot permissions")
		interactionResponseError(s, i, "Server setup not complete, please use the /status command.", true)
		return
	}
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		log.WithError(err).Error("Couldn't query channel")
		return
	}

	commandName := i.ApplicationCommandData().Name
	if handler, ok := commands.handlers[commandName]; ok {
		ctx := context.WithValue(ctx, log.Key, log.Fields{
			"author_id":    commandAuthor.ID,
			"channel_id":   i.ChannelID,
			"guild_id":     i.GuildID,
			"user":         commandAuthor.Username,
			"channel_name": channel.Name,
			"command":      commandName,
		})
		log.WithContext(ctx).Info("invoking standard command")
		handler(ctx, s, i)
		return
	}
}
