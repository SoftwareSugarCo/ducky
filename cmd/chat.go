/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with the OpenAI API",
	Long:  `This starts a dialogue with Ducky that allows ducky to remember things you tell it during the conversation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the api key
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
			return
		}

		// Get the target model
		model := GetModel(model)
		fmt.Println("Using model: " + model)

		// Initialize a scanner to read user input.
		scanner := bufio.NewScanner(os.Stdin)

		// Initialize a slice to store the conversation history.
		var conversation = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are helpful coding assistant named ducky; you always explain coding concepts in a way that is easy to understand and you always document your coding examples. You speak in a friendly tone and you are very patient but you don't mind swearing. You are a very good teacher.",
			},
		}

		var multiLine bool
		var messageLines []string

		fmt.Println("Ducky: Yes, How may I help you?")

		for {
			if !multiLine {
				fmt.Print("You: ")
			}
			// Read the user input.
			scanner.Scan()
			userInput := scanner.Text()

			if strings.ToLower(userInput) == "/done" || strings.ToLower(userInput) == "/exit" || strings.ToLower(userInput) == "/quit" || strings.ToLower(userInput) == "done" || strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" {
				// Exit the loop if the user types 'done'.
				break
			}

			if multiLine {
				if strings.EqualFold(userInput, "/end") {
					multiLine = false
					userInput = strings.Join(messageLines, "\n")
					messageLines = nil
				} else {
					messageLines = append(messageLines, userInput)
					continue
				}
			}

			if strings.EqualFold(userInput, "/multiline") || strings.EqualFold(userInput, "/ml") || strings.EqualFold(userInput, "/multi") || strings.EqualFold(userInput, "/m") {
				multiLine = true
				continue
			}

			// Append the user input to the conversation history.
			conversation = append(conversation, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})

			// Interact with the GPT API using the conversation history.
			// Replace the following line with the actual API call and response.
			gptResponse, err := chatGPT(apiKey, conversation, model)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Append the GPT response to the conversation history.
			conversation = append(conversation, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: gptResponse,
			})

			// Display the GPT response.
			fmt.Println("Ducky: " + gptResponse)

			// Attempt to extract code from the response
			codeSnips := ExtractCodeSnippets(gptResponse)
			// Check if there are any code snippets
			if len(codeSnips) > 0 {
				fmt.Println(DividerStr)
				_ = printFormattedCodeSnippets(codeSnips)
				fmt.Println(DividerStr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chatCmd.Flags().StringVarP(&model, "chat_model", "g", "gpt4", "Target model to chat with. Default is GPT4")
}

// chatGPT starts a new client to the OpenAI API and sends a query to the specified model.
func chatGPT(apiKey string, conversation []openai.ChatCompletionMessage, model string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: conversation,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
