package main

import (
	"gaspartv/GO-chatbot-com-gemini/src/configs"
	routers "gaspartv/GO-chatbot-com-gemini/src/routes"
)

var (
	logger configs.Logger
)

func main() {
	logger = *configs.GetLogger("main")

	routers.InitializeRoutes()
}
