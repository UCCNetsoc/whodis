package verify

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/models"
)

type stateMap map[string]func(m *StateParams) error

func (s stateMap) get(m *StateParams) (func(m *StateParams) error, bool) {
	if m.User == nil || m.Guild == nil {
		return nil, false
	}
	out, ok := s[m.User.ID+m.Guild.ID]
	return out, ok
}

func (s stateMap) set(m *StateParams, p func(m *StateParams) error) (string, error) {
	return models.DBClient.CreateUser(m.User.ID, m.Guild.ID)
}

var (
	state   = stateMap{}
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
	if handler, ok := state.get(s); ok {
		handler(s)
	} else if err := createLink(s); err != nil {
		return err
	}
	return nil
}

func createLink(m *StateParams) error {
	channel, err := session.UserChannelCreate(m.User.ID)
	if err != nil {
		return err
	}
	short, err := state.set(m, handleSuiteAuth)
	if err != nil {
		return err
	}
	session.ChannelMessageSend(channel.ID, fmt.Sprintf("To register for %s, click on the following link: http://%s/discord/auth?i=%s", m.Guild.Name, viper.GetString("api.hostname"), short))
	return nil
}

func handleSuiteAuth(m *StateParams) error {
	return nil
}
