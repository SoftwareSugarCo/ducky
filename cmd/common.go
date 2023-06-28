package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"io"
	"os"
	"regexp"
	"strings"
)

// GetModel returns the exact model string for the model provided. If the model is unrecognized or not passed in,
// the GPT4 model is selected by default.
func GetModel(model string) string {
	modelLower := strings.ToLower(model)

	switch modelLower {
	case openai.GPT4, GPT4:
		return openai.GPT4
	case openai.GPT3Dot5Turbo, GPT3, GPT3_5:
		return openai.GPT3Dot5Turbo
	case openai.GPT3Davinci, GPT3DAVINCI:
		return openai.GPT3Davinci
	default:
		fmt.Printf("Model unknown (%s).. defaulting to %s\n", modelLower, defaultModel)
		return openai.GPT3Dot5Turbo
	}
}

// PrintFormattedCodeSnippets takes any extracted code snippets and tries to format them based on what language it is.
func PrintFormattedCodeSnippets(snippetsByLang map[string][]string) error {
	for lang, snippets := range snippetsByLang {
		// Detect the language of the code snippet
		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		// Apply the formatting style
		style := styles.Get("solarized-dark256") // Choose a style that works well with your terminal
		if style == nil {
			style = styles.Fallback
		}

		// Format the code snippet for terminal output
		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}

		for _, snippet := range snippets {
			iterator, err := lexer.Tokenise(nil, snippet)
			if err != nil {
				return err
			}

			// Print the highlighted code snippet
			err = formatter.Format(os.Stdout, style, iterator)
			fmt.Println()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ExtractCodeSnippets parses the response and extracts the language and code snippet from a proper Markdown snippet.
// The return is a map with language as the key and list of snippets as the value(s).
//
//   - "go" => ["func main() {}", "func test() {}"]
//   - Snippets start with 3 accent symbols followed immediately by the language name and end with 3 accent symbols.
func ExtractCodeSnippets(resp string) map[string][]string {
	re := regexp.MustCompile("(?s)```(\\w+\\n)?(.+?)```(?s)")
	matches := re.FindAllStringSubmatch(resp, -1)
	snippets := make(map[string][]string)

	for _, match := range matches {
		lang := "bash"
		snippet := strings.Trim(match[2], "\n")

		if len(match[1]) > 0 {
			lang = strings.Trim(match[1], "\n")
		}
		snippets[lang] = append(snippets[lang], snippet)
	}

	return snippets
}

// SendToGPT starts a new client to the OpenAI API and sends a query to the specified model.
func SendToGPT(apiKey, model string, prompts []openai.ChatCompletionMessage) (string, error) {
	gptClient := openai.NewClient(apiKey)
	resp, err := gptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: prompts,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func SendToGPTStreamResponse(apiKey, model string, prompts []openai.ChatCompletionMessage) (string, error) {
	hiYellowPrint := color.New(color.Bold, color.FgHiYellow).PrintfFunc()
	yellowPrint := color.New(color.FgYellow).PrintfFunc()
	respBuilder := strings.Builder{}
	gptClient := openai.NewClient(apiKey)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: prompts,
		Stream:   true,
	}
	stream, err := gptClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	hiYellowPrint("Ducky: ")
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\n")
			respBuilder.Write([]byte("\n"))
			return respBuilder.String(), nil
		}

		if err != nil {
			return "", err
		}

		respBuilder.Write([]byte(res.Choices[0].Delta.Content))
		yellowPrint(res.Choices[0].Delta.Content)
	}
}
