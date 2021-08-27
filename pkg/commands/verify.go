package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"github.com/uccnetsoc/whodis/pkg/utils"
)

// VerifyCommand inits the verification process.
func VerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	guild, _ := s.Guild(i.GuildID)
	if user == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command is only available from inside a valid server",
				Flags:   1 << 6,
			},
		})
		return
	}
	uid, err := utils.Encrypt(user.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		log.Println(err)
		return
	}
	gid, err := utils.Encrypt(guild.ID, []byte(viper.GetString("api.secret")))
	if err != nil {
		log.Println(err)
		return
	}
	encoded, err := utils.Encrypt(fmt.Sprintf("%s.%s", uid, gid), []byte(viper.GetString("api.secret")))
	if err != nil {
		log.Println(err)
		return
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hey **%s**! Welcome to **%s**!", user.Username, guild.Name),
			Flags:   1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Click here to register",
							Style:    discordgo.LinkButton,
							Disabled: false,
							URL:      fmt.Sprintf("http://%s/discord/auth?state=%s", viper.GetString("api.hostname"), encoded),
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	time.AfterFunc(time.Second*15, func() {
		_, err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Button has timed out, use /verify to try again!",
							Style:    discordgo.LinkButton,
							Disabled: true,
							URL:      "https://netsoc.co/rk",
						},
					},
				},
			},
		})
	})
}
