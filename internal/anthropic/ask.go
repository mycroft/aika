package anthropic

import (
	"bytes"
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func Ask(apiKey, modelName, systemPrompt, prompt string, stream bool) (string, error) {
	maxTokens := int64(16384)
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	messages := []anthropic.MessageParam{
		anthropic.NewAssistantMessage(anthropic.NewTextBlock(systemPrompt)),
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}

	model := anthropic.ModelClaude3_7SonnetLatest

	if stream {
		var responseBuffer bytes.Buffer

		stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
			Model:     model,
			MaxTokens: maxTokens,
			Messages:  messages,
		})

		message := anthropic.Message{}
		for stream.Next() {
			event := stream.Current()
			err := message.Accumulate(event)
			if err != nil {
				panic(err)
			}

			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					responseBuffer.WriteString(deltaVariant.Text)
					fmt.Print(deltaVariant.Text)
				}

			}

			if stream.Err() != nil {
				panic(stream.Err())
			}
		}

		return responseBuffer.String(), nil
	} else {
		message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
			MaxTokens: maxTokens,
			Messages:  messages,
			Model:     model,
		})
		if err != nil {
			panic(err.Error())
		}

		return message.Content[0].Text, nil

	}
}
