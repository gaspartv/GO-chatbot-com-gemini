package routers

import (
	"gaspartv/GO-chatbot-com-gemini/src/validations"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeRoutes() {
	router := gin.Default()

	router.Static("/public", "./public")

	InitializeGeminiRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("0.0.0.0:" + validations.LoadEnv().Port)
}
