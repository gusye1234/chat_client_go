package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var openaiApiKey string = os.Getenv("OPENAI_API_KEY")
	fmt.Println("OpenAI API Key: " + openaiApiKey)

	client := openai.NewClient(openaiApiKey)
	appContext := context.Background()
	var messages []openai.ChatCompletionMessage
	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("! Error %v", err)
			continue
		}
		trimLine := strings.Trim(line, " \n")
		if trimLine == "" {
			fmt.Println("! Please input something")
			continue
		}

		var currentPack openai.ChatCompletionMessage = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: trimLine,
		}

		messages = append(messages, currentPack)

		openaiRequest := openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		}
		stream, err := client.CreateChatCompletionStream(
			appContext,
			openaiRequest,
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		var currentResponse string = ""
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Printf("\n")
				break
			}
			if err != nil {
				fmt.Printf("stream.Recv error: %v\n", err)
				break
			}
			fmt.Printf(resp.Choices[0].Delta.Content)
			currentResponse += resp.Choices[0].Delta.Content
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: currentResponse,
		})
	}
}
