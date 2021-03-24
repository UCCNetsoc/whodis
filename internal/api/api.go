package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uccnetsoc/veribot/docs"
)

func InitAPI() {
	docs.SwaggerInfo.Title = viper.GetString("api.title")
	docs.SwaggerInfo.Description = viper.GetString("api.description")
	docs.SwaggerInfo.Version = viper.GetString("api.version")
	docs.SwaggerInfo.BasePath = viper.GetString("api.path")
	docs.SwaggerInfo.Host = viper.GetString("api.hostname")

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
