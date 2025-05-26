package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gaspartv/GO-chatbot-com-gemini/src/validations"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func GeminiHandler(ctx *gin.Context) {
	var request interface{}
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"data": request,
	}

	body, err := json.Marshal(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar o body"})
		return
	}

	text, ok := bodyMap["data"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'data' não encontrado"})
		return
	}

	textStr, err := json.Marshal(text)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar o texto"})
		return
	}

	resultText := geminiAi(string(textStr))
	if resultText == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar o texto"})
		return
	}

	audio := textToAudio(resultText)
	if audio == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar o áudio"})
		return
	}

	apiUrl := validations.LoadEnv().API_URL

	ctx.JSON(http.StatusOK, gin.H{"link": apiUrl + audio})
}

func textToAudio(text string) string {
	cfg := validations.LoadEnv()
	apiKey := cfg.OpenAI_API_Key
	model := cfg.OpenAI_API_Model
	voice := cfg.OpenAI_API_Voice
	instructions := cfg.OpenAI_API_Instructions
	url := cfg.OpenAI_API_URL
	url += "audio/speech"

	data := map[string]interface{}{
		"model":        model,
		"input":        text,
		"voice":        voice,
		"instructions": instructions,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao codificar JSON:", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro na requisição:", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro na resposta:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Erro na API (%d): %s\n", resp.StatusCode, string(body))
		return ""
	}

	filename := fmt.Sprintf("public/speech_%d.mp3", os.Getpid())
	outFile, err := os.Create("./" + filename)
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return ""
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("Erro ao salvar o áudio:", err)
		return ""
	}

	return filename
}

func geminiAi(text string) string {
	ctx := context.Background()

	cfg := validations.LoadEnv()
	apiKey := cfg.GenAI_API_Key
	model := cfg.GenAI_API_Model

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	newText := "Resposta deve ser em português. O texto deve ser montado para que possa ser lido por uma IA 'text-to-speech' (TTS). O texto deve ser claro e conciso, evitando jargões técnicos. O objetivo é fornecer uma resposta útil, curta, direta e compreensível para o usuário." + text

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(newText),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return result.Text()
}
