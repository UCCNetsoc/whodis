package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (c *Client) CreateUser(discord_id string) error {
	if discord_id == "" {
		return fmt.Errorf("discord id is empty")
	}

	user := &User{DiscordID: discord_id}

	// if user not found
	err := c.conn.First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create user
			return c.conn.Create(user).Error
		}
		return err
	}
	// if user is found, return nil
	return nil
}

func (c *Client) UpdateMailDomain(discord_id, mail_domain string) {
	user := &User{DiscordID: discord_id}
	c.conn.First(user)
	user.MailDomain = mail_domain
	c.conn.Save(user)
}
