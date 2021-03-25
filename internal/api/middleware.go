package api

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Request.Cookie("discord_id")
		if err != nil {
			c.Redirect(http.StatusFound, "/discord/login")
			return
		}
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			return viper.GetString("api.secret"), nil
		})
		if err != nil {
			log.Println("Couldn't parse token")
			c.Status(http.StatusBadRequest)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println(claims)
			c.Next()
		}
	}
}
