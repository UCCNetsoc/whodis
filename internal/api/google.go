package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/uccnetsoc/whodis/pkg/models"

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
			// scope to view email address
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func googleCheckLogin(state string) string {
	return googleConf.AuthCodeURL(state)
}

func googleLoginHandler(c *gin.Context) {
	log.Println(googleConf)
	state := ``
	c.Redirect(http.StatusTemporaryRedirect, googleConf.AuthCodeURL(googleCheckLogin(state)))
}

func googleAuthHandler(c *gin.Context) {
	token, err := googleConf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := googleConf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading json body\n", err)
		c.Status(http.StatusBadRequest)
		return
	}

	jsonResp := make(map[string]interface{})

	if err = json.Unmarshal(body, &jsonResp); err != nil {
		log.Println("Couldn't unmarshall json data\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	tokenString, err := c.Request.Cookie("discord_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/discord/login")
		return
	}

	jwtToken, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("api.secret")), nil
	})

	if err != nil {
		log.Println("Couldn't parse token\n", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		models.DBClient.UpdateMailDomain(claims["ID"].(string), jsonResp["hd"].(string))
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
