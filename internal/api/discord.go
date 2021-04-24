package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/uccnetsoc/whodis/pkg/models"
)

var discordConf *oauth2.Config

// InitOAuth sets up OAuth client
func InitDiscordOAuth() {
	discordConf = &oauth2.Config{
		ClientID:     viper.GetString("oauth.discord.id"),
		ClientSecret: viper.GetString("oauth.discord.secret"),
		RedirectURL:  "http://" + viper.GetString("api.hostname") + ":8080" + "/discord/auth",
		Scopes: []string{
			"identify",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func discordCheckLogin(state string) string {
	return discordConf.AuthCodeURL(state)
}

// func randToken() string {
// 	b := make([]byte, 32)
// 	rand.Read(b)
// 	return base64.StdEncoding.EncodeToString(b)
// }

func discordLoginHandler(c *gin.Context) {
	state := ``
	c.Redirect(http.StatusTemporaryRedirect, discordCheckLogin(state))
}

func discordAuthHandler(c *gin.Context) {
	token, err := discordConf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := discordConf.Client(context.Background(), token)
	resp, err := client.Get("https://discordapp.com/api/users/@me")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	var discordResp *models.DiscordResp
	json.Unmarshal(body, &discordResp)

	log.Println(discordResp.ID)

	// create JWT
	expirationTime := time.Now().Add(time.Hour * 24 * 30).Unix()
	claims := &models.JWT{
		ID: discordResp.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
		},
	}

	client_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := client_token.SignedString([]byte(viper.GetString("api.secret")))
	if err != nil {
		log.Println("Couldn't generate api token")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	params, ok := c.Get("i")
	if !ok {
		log.Println("No params")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	short, ok := params.(string)
	if !ok {
		log.Println("Param not string")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	guild, err := models.DBClient.GetGuildFromShort(short)
	if err != nil {
		log.Println("Guild does not exist")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("discord_id", tokenString, 0, "/", viper.GetString("api.hostname"), false, true)
	c.SetCookie("discord_guild_id", guild.ID, 0, "/", viper.GetString("api.hostname"), false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/google/login")
}
