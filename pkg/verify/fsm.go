package verify

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type stateMap map[string]func(m *StateParams) error

func (s stateMap) get(m *StateParams) (func(m *StateParams) error, bool) {
	if m.User == nil || m.Guild == nil {
		return nil, false
	}
	out, ok := s[m.User.ID+m.Guild.ID]
	return out, ok
}

func (s stateMap) set(m *StateParams, p func(m *StateParams) error) {
	s[m.User.ID+m.Guild.ID] = p
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
	claims := jwt.MapClaims{
		"user_id":  m.User.ID,
		"guild_id": m.Guild.ID,
	}
	client_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := client_token.SignedString([]byte(viper.GetString("api.secret")))
	session.ChannelMessageSend(channel.ID, fmt.Sprintf("To register for %s, click on the following link: http://%s/discord/auth?i=%s", m.Guild.Name, viper.GetString("api.hostname"), tokenString))
	state.set(m, handleSuiteAuth)
	return nil
}

func handleSuiteAuth(m *StateParams) error {
	return nil
}
