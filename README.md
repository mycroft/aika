# Aika

Aika is a command-line tool that leverages AI to assist with various tasks such as code review, commit message generation, and README generation. It supports multiple AI providers, including OpenAI, Mistral, and Anthropic.

## Features

- **AI Providers**: Choose between OpenAI, Mistral, or Anthropic as your AI provider.
- **Model Selection**: Specify the model name for your chosen AI provider.
- **Streaming**: Stream the AI response for real-time feedback.
- **Diff Input**: Use the output of `git diff` or `git diff --cached` as input.
- **Tree Input**: Use the entire current tree as input.
- **Feature Selection**: Specify the feature to use, such as code review, commit message generation, or README generation.

## Installation

To install Aika, you can use the following command:

```sh
go get -u github.com/mycroft/aika
```

## Usage

Aika can be used with various options to customize its behavior. Here are some examples:

### Basic Usage

```sh
aika -ai mistral -model <model-name> -feature code-review
```

### Using Git Diff as Input

```sh
aika -ai openai -model <model-name> -feature commit-message -diff
```

### Using Git Diff --Cached as Input

```sh
aika -ai anthropic -model <model-name> -feature readme -diff-cached
```

### Using the Entire Tree as Input

```sh
aika -ai mistral -model <model-name> -feature code-review -tree
```

## Options

- `-ai`: AI provider (openai, mistral, or anthropic).
- `-model`: Model name (default depends on the provider).
- `-stream`: Stream the response.
- `-diff`: Use `git diff` as input.
- `-diff-cached`: Use `git diff --cached` as input.
- `-tree`: Use the entire current tree as input.
- `-feature`: Feature to use (code-review, commit-message, or readme).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.