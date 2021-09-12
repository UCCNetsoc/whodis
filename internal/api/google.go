package api

import (
	"net/http"

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
		RedirectURL:  viper.GetString("api.url") + "/google/auth",
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
	if c.Query("hd") != viper.GetString("oauth.google.domain") {
		resultTemplate.Execute(c.Writer, AccessErrorResponse(http.StatusBadRequest, "Invalid oauth domain. Wanted: "+viper.GetString("oauth.google.domain"), nil))
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/verify?state="+c.Query("state"))
}
