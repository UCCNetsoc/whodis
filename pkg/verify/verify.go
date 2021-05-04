package verify

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/models"
	"gorm.io/gorm"
)

func getState(m *StateParams) (guild *models.Guild, found bool, err error) {
	if m.User == nil || m.Guild == nil {
		err = errors.New("No user or guild set.")
		return
	}
	guild, err = models.DBClient.GetGuildFromID(m.User.ID, m.Guild.ID)
	found = !errors.Is(err, gorm.ErrRecordNotFound)
	return
}

func setState(m *StateParams, p func(m *StateParams) error) (string, error) {
	return models.DBClient.CreateUser(m.User.ID, m.Guild.ID)
}

var (
	session *discordgo.Session
)

func Init(s *discordgo.Session) {
	session = s
}

type StateParams struct {
	User  *discordgo.User
	Guild *discordgo.Guild
}

// Transition from one state to another base on the decision tree.
func Transition(s *StateParams) error {
	guild, found, err := getState(s)
	switch {
	case !found:
		return createLink(s)
	case err != nil:
		return err
	case guild.Verified:
		return handleSuiteAuth(s)
	}
	return nil
}

func createLink(m *StateParams) error {
	channel, err := session.UserChannelCreate(m.User.ID)
	if err != nil {
		return err
	}
	short, err := setState(m, handleSuiteAuth)
	if err != nil {
		return err
	}
	session.ChannelMessageSend(channel.ID, fmt.Sprintf("To register for %s, click on the following link: http://%s/discord/auth?i=%s", m.Guild.Name, viper.GetString("api.hostname"), short))
	return nil
}

func handleSuiteAuth(m *StateParams) error {
	return nil
}
