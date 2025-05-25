package handlers

import "gaspartv/GO-chatbot-com-gemini/src/configs"

var (
	logger *configs.Logger
)

func InitializeHandlers() {
	logger = configs.GetLogger("handlers")
}
