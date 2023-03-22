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

var output string

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Query the OpenAI API with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		apiKey := viper.GetString("api_key")

		if apiKey == "" {
			fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
			return
		}

		resp, err := queryChatGPT(apiKey, prompt)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Switch statement for the output flag
		switch output {
		case "stdout":
			fmt.Println(resp)
		case "json":
			fmt.Println(resp)
		case "file":
			fmt.Println(resp)
		default:
			fmt.Println(resp)

		}

		codeSnippets := extractCodeSnippets(resp)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("##################################################")
			fmt.Println("Code snippets from response: ")
			_ = printFormattedCodeSnippets(codeSnippets)
			fmt.Println("##################################################")
		}

	},
}

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
	rootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// askCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// askCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	askCmd.Flags().StringVarP(&output, "output", "o", "stdout", "Output format. Options are stdout, json, and file")
}

func queryChatGPT(apiKey string, prompt string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
