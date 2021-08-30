package api

import (
	"fmt"
	"net/http"

	"github.com/Strum355/log"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConf *oauth2.Config

// InitOAuth sets up OAuth client
func InitGoogleOAuth() {
	googleConf = &oauth2.Config{
		ClientID:     viper.GetString("oauth.google.id"),
		ClientSecret: viper.GetString("oauth.google.secret"),
		RedirectURL:  "http://" + viper.GetString("api.hostname") + ":8080" + "/google/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func googleLoginHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, googleConf.AuthCodeURL(c.Query("state")))
}

func googleAuthHandler(c *gin.Context) {
	if c.Query("hd") != "umail.ucc.ie" {
		log.Error("Invalid umail address")
		c.Writer.WriteString("Invalid umail address")
		c.Status(http.StatusNotAcceptable)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/verify?state=%s", c.Query("state")))
}
