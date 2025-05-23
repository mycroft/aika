package mistral

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const APIURL = "https://api.mistral.ai/v1/chat/completions"

type MistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralRequest struct {
	Model    string           `json:"model"`
	Messages []MistralMessage `json:"messages"`
	Stream   bool             `json:"stream,omitempty"`
}

type MistralResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type MistralStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
		Logprobs     *struct {
			TokenLogprobs []float64            `json:"token_logprobs"`
			Tokens        []string             `json:"tokens"`
			TopLogprobs   []map[string]float64 `json:"top_logprobs"`
			TextOffset    []int                `json:"text_offset"`
		} `json:"logprobs"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

func Ask(apiKey, model, systemPrompt, prompt string, stream bool) (string, error) {
	reqBody := MistralRequest{
		Model: model,
		Messages: []MistralMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Stream: stream,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", APIURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if stream {
		var responseBuffer bytes.Buffer
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", fmt.Errorf("error reading response: %w", err)
			}

			// Skip empty lines or lines that don't start with data:
			if len(line) == 0 || line == "\n" {
				continue
			}

			// If the API streams with a prefix (e.g., "data: "), strip it:
			const prefix = "data: "
			if len(line) > len(prefix) && line[:len(prefix)] == prefix {
				line = line[len(prefix):]
			}

			var mistralResp MistralStreamResponse
			if err := json.Unmarshal([]byte(line), &mistralResp); err != nil {
				continue // skip malformed lines
			}

			// fmt.Println(mistralResp)
			if len(mistralResp.Choices) == 1 {
				fmt.Print(mistralResp.Choices[0].Delta.Content)
				responseBuffer.WriteString(mistralResp.Choices[0].Delta.Content)
			}
		}

		return responseBuffer.String(), nil
	} else {
		var mistralResp MistralResponse
		if err := json.NewDecoder(resp.Body).Decode(&mistralResp); err != nil {
			return "", fmt.Errorf("error decoding response: %w", err)
		}

		if len(mistralResp.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		return mistralResp.Choices[0].Message.Content, nil
	}
}
