package models

import (
	"errors"

	"gorm.io/gorm"
)

type Guild struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Verified   bool   `json:"verified"`
	MailDomain string `json:"mail_domains,omitempty"`
}

type User struct {
	DiscordID string   `gorm:"primaryKey" json:"discord_id,omitempty"`
	Guilds    []*Guild `gorm:"foreignKey:ID" json:"guild_id,omitempty"`
}

func (c *Client) CreateUser(discord_id string, guild_id string) error {
	user := &User{DiscordID: discord_id}
	err := c.conn.First(user).Preload("Guilds").Error
	guidExists := false
	for _, g := range user.Guilds {
		if g.ID == guild_id {
			guidExists = true
		}
	}
	if !guidExists {
		user.Guilds = append(user.Guilds, &Guild{ID: guild_id})
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.conn.Create(user).Error
		}
		return err
	}
	return c.conn.Save(user).Error
}
