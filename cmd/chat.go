/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"ducky/util"
	"fmt"
	"github.com/fatih/color"
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
	hiYellowPrint := color.New(color.Bold, color.FgHiYellow).PrintfFunc()
	yellowPrint := color.New(color.Bold, color.FgYellow).PrintfFunc()
	yellowStr := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	redStr := color.New(color.FgHiRed).SprintFunc()
	greenStr := color.New(color.FgHiGreen).SprintFunc()
	cyanStr := color.New(color.Bold, color.FgHiCyan).SprintFunc()
	hiBluePrint := color.New(color.Bold, color.FgHiBlue).PrintfFunc()
	bluePrint := color.New(color.FgBlue).PrintfFunc()
	// Get the api key
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
		return
	}

	var streamModeStr string
	if StreamMode {
		streamModeStr = greenStr(StreamMode)
	} else {
		streamModeStr = redStr(StreamMode)
	}

	var toFileStr string
	if ToFile {
		toFileStr = greenStr(ToFile)
	} else {
		toFileStr = redStr(ToFile)
	}

	util.PrintDuckyHeader()

	// Get the target model
	model := GetModel(model)

	util.PrintBox("SETTINGS", map[string]string{"Model": cyanStr(model), "Mode": yellowStr("Chat"), "Stream": streamModeStr, "ToFile": toFileStr})
	util.PrintBox("COMMANDS", map[string]string{"Exit": "/exit, /quit, /q, /done", "Multi-line": "/m, /ml, /multi, /multiline", "End multi-line": "/end", "Set Personality": "/whoami"})

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

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	hiYellowPrint("Ducky: ")
	yellowPrint("Yes, How may I help you?\n")
	for {
		if !multiLine {
			hiBluePrint("You: ")
		}
		// Read the user input.
		scanner.Scan()
		userInput := scanner.Text()

		// Before sending the query to OpenAI, check user input for any commands
		switch strings.ToLower(userInput) {
		case "/q", "/quit", "/exit", "/done": // Exit commands
			stopChat = true
			hiYellowPrint("\nDucky: ")
			hiYellowPrint("Goodbye!\n")
			break
		case "/m", "/ml", "/multi", "/multiline": // Multiline commands
			multiLine = true
			continue
		case "/whoami": // Set system personality command
			setPersonality = true
			yellowPrint("Who am I?\n")
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

			yellowPrint("That's who I'll be then...\n")
			continue
		}

		// Append the user input to the conversation history.
		conversation = append(conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		var (
			gptResponse string
			err         error
		)
		if StreamMode {
			gptResponse, err = SendToGPTStreamResponse(apiKey, model, conversation)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			gptResponse, err = SendToGPT(apiKey, model, conversation)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Display the GPT response.
			hiYellowPrint("\nDucky: ")
			yellowPrint("%s\n\n", gptResponse)
			bluePrint("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
		}

		// Append the GPT response to the conversation history.
		conversation = append(conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: gptResponse,
		})

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

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&model, "chat_model", "c", "gpt3", "Target model to chat with. Default is GPT3")
}
