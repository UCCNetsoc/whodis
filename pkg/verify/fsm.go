package verify

import "github.com/bwmarrin/discordgo"

var state = map[string]func() error{}

// Transition from one state to another base on the decision tree.
func Transition(m *discordgo.Message) error {
	switch m.Content {
	default:
		if err := createLink(); err != nil {
			return err
		}
	}
	return nil
}

func createLink() error {
	return nil
}
