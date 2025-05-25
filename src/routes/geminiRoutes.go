package routers

import (
	"gaspartv/GO-chatbot-com-gemini/src/handlers"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/override/docs"
)

func InitializeGeminiRoutes(router *gin.Engine) {
	handlers.InitializeHandlers()

	baseUrl := "/api/v1/gemini"

	docs.SwaggerInfo.BasePath = baseUrl

	v1 := router.Group(baseUrl)
	{
		v1.POST("/", handlers.GeminiHandler)
	}
}
