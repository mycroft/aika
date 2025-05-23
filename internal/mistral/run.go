package mistral

import (
	"fmt"
	"os"
)

// Possible models:
// - devstral-small-2505
// - mistral-small-2503

const APIKeyEnv = "MISTRAL_API_KEY"

func Run(defaultPrompt, modelName, input string, stream bool) string {
	apiKey := os.Getenv(APIKeyEnv)
	if apiKey == "" {
		fmt.Println("Please set the MISTRAL_API_KEY environment variable.")
		os.Exit(1)
	}

	response, err := Ask(apiKey, modelName, defaultPrompt, string(input), stream)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return response
}
