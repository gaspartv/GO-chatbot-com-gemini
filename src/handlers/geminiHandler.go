package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Defina a variável de ambiente API_URL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"link": apiUrl + audio})
}

func textToAudio(text string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Defina a variável de ambiente OPENAI_API_KEY")
		return ""
	}

	url := os.Getenv("OPENAI_API_URL")
	if url == "" {
		fmt.Println("Defina a variável de ambiente OPENAI_API_URL")
		return ""
	}
	url += "audio/speech"

	model := os.Getenv("OPENAI_API_MODEL")
	if model == "" {
		fmt.Println("Defina a variável de ambiente OPENAI_API_MODEL")
		return ""
	}

	voice := os.Getenv("OPENAI_API_VOICE")
	if voice == "" {
		fmt.Println("Defina a variável de ambiente OPENAI_API_VOICE")
		return ""
	}

	instructions := os.Getenv("OPENAI_API_INSTRUCTIONS")
	if instructions == "" {
		fmt.Println("Defina a variável de ambiente OPENAI_API_INSTRUCTIONS")
		return ""
	}

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

	apiKey := os.Getenv("GENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Defina a variável de ambiente GENAI_API_KEY")
	}

	model := os.Getenv("GENAI_API_MODEL")
	if model == "" {
		log.Fatal("Defina a variável de ambiente GENAI_API_MODEL")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(text),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return result.Text()
}
