package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	gorm.Model
	GuildID        string `gorm:"primaryKey"`
	Roles          []string
	WelcomeChannel string
	Domains        []ConfigItem
}

type ConfigItem struct {
	gorm.Model
	value string
}

func (c *Client) SetConfigItem(guid string, key string, value interface{}) error {
	return c.conn.Transaction(func(tx *gorm.DB) error {
		key = strings.ReplaceAll(strings.ToLower(key), "_", "")
		c := &Config{GuildID: guid}
		if err := tx.First(c).Preload(clause.Associations).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(c).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		refConfig := reflect.ValueOf(c).Elem()
		refField := refConfig.FieldByName(key)
		if refField.Kind() == reflect.Slice {
			if refValue := reflect.ValueOf(value); refValue.Kind() != reflect.Slice {
				return fmt.Errorf("value must be a slice but got %s", refValue.Kind().String())
			}
			valueStr, ok := value.([]string)
			if !ok {
				return fmt.Errorf("invalid value type")
			}
			items := []ConfigItem{}
			for _, v := range valueStr {
				items = append(items, ConfigItem{value: v})
			}
			return tx.Set(key, items).Error
		} else {
			return tx.Set(key, value).Error
		}
	})
}
