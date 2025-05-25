package main

import (
	"fmt"
	"gaspartv/GO-chatbot-com-gemini/src/configs"
	routers "gaspartv/GO-chatbot-com-gemini/src/routes"
	"os"

	"github.com/joho/godotenv"
)

var (
	logger configs.Logger
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		os.Exit(1)
	}

	logger = *configs.GetLogger("main")

	routers.InitializeRoutes()
}
