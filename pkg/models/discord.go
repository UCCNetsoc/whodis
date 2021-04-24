package models

import (
	"errors"
	"math/rand"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Guild struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Verified   bool   `json:"verified"`
	MailDomain string `json:"mail_domains,omitempty"`
	Short      string `gorm:"unique" json:"short"`
}

type User struct {
	DiscordID string   `gorm:"primaryKey" json:"discord_id,omitempty"`
	Guilds    []*Guild `gorm:"foreignKey:ID" json:"guild_id,omitempty"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (c *Client) CreateUser(discord_id string, guild_id string) (string, error) {
	user := &User{DiscordID: discord_id}
	var short string
	err := c.conn.Transaction(func(tx *gorm.DB) error {
		short := randStringRunes(viper.GetInt("api.sluglength"))
		if err := tx.Find(&Guild{}, "short = ?", short).Error; err != nil {
			return err
		}
		err := c.conn.First(user).Preload("Guilds").Error
		guidExists := false
		for _, g := range user.Guilds {
			if g.ID == guild_id {
				guidExists = true
			}
		}
		if !guidExists {
			user.Guilds = append(user.Guilds, &Guild{ID: guild_id, Short: short})
		}
		if err != nil {
			return err
		}
		return tx.Save(user).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.CreateUser(discord_id, guild_id)
	}
	return short, err
}
