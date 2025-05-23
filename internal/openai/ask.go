package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Ask(apiKey, model, systemPrompt, prompt string, stream bool) (string, error) {
	client := openai.NewClient(apiKey)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	if stream {
		resp, err := client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			return "", fmt.Errorf("error creating chat completion stream: %w", err)
		}
		defer resp.Close()

		var responseBuffer bytes.Buffer

		for {
			content, err := resp.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("error while streaming response...")
				os.Exit(1)
			}

			responseBuffer.WriteString(content.Choices[0].Delta.Content)
			fmt.Print(content.Choices[0].Delta.Content)
		}

		return responseBuffer.String(), nil
	} else {
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			os.Exit(1)
		}

		return resp.Choices[0].Message.Content, nil
	}
}
