/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const DividerStr = "##################################################"

var (
	output, model string
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [prompt]",
	Short: "Query the OpenAI API with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		apiKey := viper.GetString("api_key")

		if apiKey == "" {
			fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
			return
		}

		// Default model selected is GPT3.5 Turbo
		model = getModel(model)

		// Send out the query to OpenAI
		resp, err := queryChatGPT(apiKey, prompt, model)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Extract any code snippets from the response and store them
		codeSnippets := extractCodeSnippets(resp)

		// Switch statement for the output flag
		switch output {
		case "json":
			// TODO: Need to implement
			fmt.Println(resp)
		case "file":
			// TODO: Need to implement
			fmt.Println(resp)
		default:
			// Default output is to stdout
			fmt.Println(resp)
			fmt.Println(DividerStr)
			_ = printFormattedCodeSnippets(codeSnippets)
			fmt.Println(DividerStr)
		}
	},
}

// getModel returns the exact model string for the model provided. If the model is unrecognized or not passed in,
// the GPT4 model is selected by default.
func getModel(model string) string {
	modelLower := strings.ToLower(model)

	switch modelLower {
	case "gpt3", "gpt3.5", "gpt3.5turbo":
		return openai.GPT3Dot5Turbo
	case "davinci", "gpt3davinci":
		return openai.GPT3Davinci
	default:
		return openai.GPT4
	}
}

// printFormattedCodeSnippets takes any extracted code snippets and tries to format them based on what language it is.
func printFormattedCodeSnippets(snippetsByLang map[string][]string) error {
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

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	queryCmd.Flags().StringVarP(&output, "output", "o", "stdout", "Output format. Options are stdout, json, and file")
	queryCmd.Flags().StringVarP(&model, "model", "m", "gpt3.5turbo", "Target model to send query to. Default is GPT3.5Turbo")
}

// queryChatGPT starts a new client to the OpenAI API and sends a query to the specified model.
func queryChatGPT(apiKey string, prompt string, model string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// extractCodeSnippets parses the response and extracts the language and code snippet from a proper Markdown snippet.
// The return is a map with language as the key and list of snippets as the value(s).
//
//   - "go" => ["func main() {}", "func test() {}"]
//   - Snippets start with 3 accent symbols followed immediately by the language name and end with 3 accent symbols.
func extractCodeSnippets(resp string) map[string][]string {
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
