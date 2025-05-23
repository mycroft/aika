package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mycroft/aika/internal/anthropic"
	"github.com/mycroft/aika/internal/mistral"
	"github.com/mycroft/aika/internal/openai"
)

func main() {
	aiKind := flag.String("ai", "mistral", "AI provider: openai, mistral, or anthropic")
	modelName := flag.String("model", "", "Model name (default depends on provider)")
	stream := flag.Bool("stream", false, "Stream response")
	diff := flag.Bool("diff", false, "Use 'git diff' as input")
	diffCached := flag.Bool("diff-cached", false, "Use 'git diff --cached' as input")
	feature := flag.String("feature", "", "Feature name to use: code-review, commit-message, or readme")

	flag.Parse()

	defaultPrompt := `
		You are a code review assistant.

		Can you review the following content, and only return possible code modification.
		Have a focus on readability, maintainability, correctness, completeness, performance,
		and security. When giving suggestions, please write each the modification suggestion
		using unified diffs in different sections with a short explaination. Please sum up
		suggestions that does not have any code modification (or only comment) to suggest
		in a dedicated section.
	`
	prompt := defaultPrompt

	switch *feature {
	case "code-review":
		prompt = defaultPrompt
	case "commit-message":
		prompt = `
		You are a commit message assistant.
		Can you generate a commit message for the following content, and only return the commit message.
	`
	case "readme":
		prompt = `
		You are a README generator.
		Can you generate a README for the following content, and only return the README, in markdown format.
	`
	}

	var input []byte
	var err error

	if *diff {
		cmd := exec.Command("git", "diff")
		input, err = cmd.Output()
	} else if *diffCached {
		cmd := exec.Command("git", "diff", "--cached")
		input, err = cmd.Output()
	} else {
		input, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	if len(input) == 0 {
		fmt.Println("Error: Input is empty.")
		return
	}

	var response string

	switch *aiKind {
	case "mistral":
		if *modelName == "" {
			*modelName = "devstral-small-2505"
		}
		response = mistral.Run(prompt, *modelName, string(input), *stream)
	case "openai":
		if *modelName == "" {
			*modelName = "gpt-4.1"
		}
		response = openai.Run(prompt, *modelName, string(input), *stream)
	case "anthropic":
		if *modelName == "" {
			*modelName = "claude-sonnet-2.7"
		}
		response = anthropic.Run(prompt, *modelName, string(input), *stream)
	default:
		fmt.Println("Unsupported AI provider. Please use 'mistral', 'openai' or 'anthropic'.")
		return
	}

	if !*stream {
		fmt.Println(response)
	}
}
