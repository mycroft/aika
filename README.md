# aika: AI Code Reviewer agent

## Build

```
$ go get -v ./...
$ go build -o . ./...
```

## Run

```sh
$ ./aika -help
Usage of ./aika:
  -ai string
        AI provider: openai, mistral, or anthropic (default "mistral")
  -model string
        Model name (default depends on provider)
  -stream
        Stream response
```

Real live example:

```sh
$ ./aika -stream < cmd/aika/main.go
Here are the suggested modifications and comments for the provided code:

### Code Modifications
...

```