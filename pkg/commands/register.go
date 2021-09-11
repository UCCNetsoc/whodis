package commands

import (
	"context"
	"errors"

	"github.com/Strum355/log"
	"github.com/spf13/viper"

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
	s.AddHandler(func(s *discordgo.Session, a *discordgo.GuildMemberAdd) {
		g, err := s.Guild(a.GuildID)
		if err != nil {
			log.WithError(err).Error("Failed to get guild")
			return
		}
		s.ChannelMessageSend(g.SystemChannelID, "Welcome to **"+g.Name+"** "+a.User.Mention()+"! Type `/verify` to get access to the full server.")
	})
}

type CommandHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) *interactionError

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
		if _, err := s.ApplicationCommandCreate(viper.GetString("discord.app.id"), "", comm); err != nil {
			log.WithError(err).Error("Failed to create command")
			return err
		}
	}
	return nil
}

func checkDirectMessage(i *discordgo.InteractionCreate) (*discordgo.User, *interactionError) {
	if i.GuildID == "" {
		return nil, &interactionError{
			errors.New("Command invoked outside of valid guild"),
			"This command is only available from inside a valid server",
		}
	}
	return i.Member.User, nil
}

// Call command handler.
func callHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var iError *interactionError
	ctx := context.Background()
	commandAuthor, iError := checkDirectMessage(i)
	if iError != nil {
		iError.Handle(s, i)
		return
	}
	commandName := i.ApplicationCommandData().Name
	roleExists, roleAccess := botPermissionCheck(s, i.GuildID)
	if commandName != "status" && !(roleExists && roleAccess) {
		iError = &interactionError{errors.New("Setup is not complete"), "Server setup not complete, please use the /status command"}
		iError.Handle(s, i)
		return
	}
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		iError = &interactionError{err, "Couldn't query channel"}
		iError.Handle(s, i)
		return
	}
	if handler, ok := commands.handlers[commandName]; ok {
		ctx := context.WithValue(ctx, log.Key, log.Fields{
			"author_id":    commandAuthor.ID,
			"channel_id":   i.ChannelID,
			"guild_id":     i.GuildID,
			"user":         commandAuthor.Username,
			"channel_name": channel.Name,
			"command":      commandName,
		})
		log.WithContext(ctx).Info("Invoking standard command")
		iError = handler(ctx, s, i)
		if iError != nil {
			iError.Handle(s, i)
		}
	}
}
