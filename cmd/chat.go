/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"ducky/util"
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
	Run:   handleChatCmd,
}

func handleChatCmd(cmd *cobra.Command, args []string) {
	// Get the api key
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
		return
	}

	util.PrintDuckyHeader()

	// Get the target model
	model := GetModel(model)

	printBox("SETTINGS", map[string]string{"Model": model, "Mode": "Chat"})
	printBox("COMMANDS", map[string]string{"Exit": "/exit, /quit, /q, /done", "Multi-line": "/m, /ml, /multi, /multiline", "End multi-line": "/end"})

	if ToFile {
		fmt.Println("ToFile mode not yet implemented")
	}

	// Initialize a scanner to read user input.
	scanner := bufio.NewScanner(os.Stdin)

	// Initialize a slice to store the conversation history.
	var conversation = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are helpful coding assistant named ducky; you always explain coding concepts in a way that is easy to understand and you always document your coding examples. You speak in a friendly tone and you are very patient but you don't mind swearing. You are a very good teacher.",
		},
	}

	var (
		multiLine      bool
		setPersonality bool
		messageLines   []string
		stopChat       bool
	)

	fmt.Println("Ducky: Yes, How may I help you?")

	for {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		if !multiLine {
			fmt.Print("You: ")
		}
		// Read the user input.
		scanner.Scan()
		userInput := scanner.Text()

		// Before sending the query to OpenAI, check user input for any commands
		switch strings.ToLower(userInput) {
		case "/q", "/quit", "/exit", "/done": // Exit commands
			stopChat = true
			fmt.Println("Ducky: Goodbye!")
			break
		case "/m", "/ml", "/multi", "/multiline": // Multiline commands
			multiLine = true
			continue
		case "/whoami": // Set system personality command
			setPersonality = true
			fmt.Println("Who am I?")
		}

		// Break from the chat loop if user had send an exit command
		if stopChat {
			break
		}

		// If multiline mode, check for the /end command to break out of it otherwise append line to message
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

		// If SetPersonality mode, user input is describing how the AI should act
		if setPersonality {
			scanner.Scan()
			personality := scanner.Text()
			setPersonality = false

			conversation = append(conversation, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: personality,
			})

			fmt.Println("That's who I'll be then...")
			continue
		}

		// Append the user input to the conversation history.
		conversation = append(conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// Interact with the GPT API using the conversation history.
		// Replace the following line with the actual API call and response.
		gptResponse, err := SendToGPT(apiKey, model, conversation)
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
		fmt.Printf("Ducky: %s\n\n", gptResponse)
		fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

		// Attempt to extract code from the response
		codeSnips := ExtractCodeSnippets(gptResponse)
		// Check if there are any code snippets
		if len(codeSnips) > 0 {
			fmt.Println(TITLEDIVSTR)
			_ = PrintFormattedCodeSnippets(codeSnips)
			fmt.Println(DIVSTR)
		}
	}
}

func printBox(label string, items map[string]string) {
	var maxLength int

	// Find the length of the longest setting key or value
	for key, value := range items {
		if len(key) > maxLength {
			maxLength = len(key)
		}

		if len(value) > maxLength {
			maxLength = len(value)
		}
	}

	// Add two to the maxLength to account for the padding on either side of the value
	maxLength += 2

	// Print the box
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Print the label line
	fmt.Printf("| %s%-*s |\n", label, ((2*maxLength)-len(label))+3, "")

	// Print the separator line
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	// Print the items
	for key, value := range items {
		fmt.Printf("| %-*s | %-*s |\n", maxLength, key, maxLength, value)
	}

	// Print the bottom of the box
	fmt.Print("+")
	for i := 0; i < (2*maxLength)+5; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	fmt.Printf("\n")
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&model, "chat_model", "c", "gpt3", "Target model to chat with. Default is GPT3")
}
