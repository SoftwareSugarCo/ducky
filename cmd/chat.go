package cmd

import (
	"bufio"
	"ducky/util"
	"fmt"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	hiYellowPrint = color.New(color.Bold, color.FgHiYellow).PrintfFunc()
	yellowPrint   = color.New(color.Bold, color.FgYellow).PrintfFunc()
	yellowStr     = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	redStr        = color.New(color.FgHiRed).SprintFunc()
	greenStr      = color.New(color.FgHiGreen).SprintFunc()
	cyanStr       = color.New(color.Bold, color.FgHiCyan).SprintFunc()
	hiBluePrint   = color.New(color.Bold, color.FgHiBlue).PrintfFunc()
	bluePrint     = color.New(color.FgBlue).PrintfFunc()
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with the OpenAI API",
	Long:  `This starts a dialogue with Ducky that allows ducky to remember things you tell it during the conversation.`,
	Run:   handleChatCmd,
}

func handleChatCmd(cmd *cobra.Command, args []string) {
	// Get the api key from environment variable
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
		return
	}

	// Determine if stream mode is to be enabled
	var streamModeStr string
	if StreamMode {
		streamModeStr = greenStr(StreamMode)
	} else {
		streamModeStr = redStr(StreamMode)
	}

	// Determine if Code snippets should be exported to a file
	var toFileStr string
	if ToFile {
		toFileStr = greenStr(ToFile)
	} else {
		toFileStr = redStr(ToFile)
	}

	// Print the DUCKY header
	util.PrintDuckyHeader()

	// Get the target model
	model := GetModel(model)

	// Print the configured settings
	util.PrintBox("SETTINGS", map[string]string{"Model": cyanStr(model), "Mode": yellowStr("Chat"), "Stream": streamModeStr, "ToFile": toFileStr})
	util.PrintBox("COMMANDS", map[string]string{"Exit": "/exit, /quit, /q, /done", "Multi-line": "/m, /ml, /multi, /multiline", "End multi-line": "/end", "Set Personality": "/whoami"})

	if ToFile {
		fmt.Println("ToFile mode not yet implemented")
	}

	// Initialize a scanner to read user input.
	scanner := bufio.NewScanner(os.Stdin)

	// Chat setting variables
	var (
		multiLine      bool
		setPersonality bool
		messageLines   []string
		stopChat       bool
	)

	var conversation []openai.ChatCompletionMessage

	// Determine Ducky's role (personality)
	getDuckySystemRole(conversation, scanner)

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

			if ToFile {
				exportSnippetsToFile(userInput, codeSnips)
			}
		}
	}
}

func exportSnippetsToFile(query string, snips map[string][]string) {
	// TODO: Need another setting for output. Default will be to a Markdown file but other iterations could split by language type
	// Check if a "ducky.md" file exists, if not create it
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory: ", err)
		return
	}
	filePath := filepath.Join(curDir, "ducky.md")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer f.Close()

	sb := strings.Builder{}
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	mdSectionHeader := fmt.Sprintf("## %s - Query\n\n", timeStamp)
	sb.WriteString(mdSectionHeader)

	sb.WriteString(query + "\n")

	for lang, snip := range snips {
		openingCodeBlock := fmt.Sprintf("```%s\n", lang)
		sb.WriteString(openingCodeBlock)
		snipStr := strings.Join(snip, "\n")
		sb.WriteString(snipStr + "\n")
		sb.WriteString("```\n\n")
	}

	_, err = f.WriteString(sb.String())
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}
	// Make sure we're always adding to file, never overwrite
}

// getDuckySystemRole requests the user to describe the system's role or defaults to a specific role
func getDuckySystemRole(conv []openai.ChatCompletionMessage, scanner *bufio.Scanner) {
	hiYellowPrint("Ducky: ")
	yellowPrint("Who am I? My default role is of a coding professor but I can try to be anything you prefer.\n")
	scanner.Scan()
	duckyRole := scanner.Text()
	if duckyRole == "" {
		duckyRole = "You are an expert coding professor. You explain concepts in great detail with simple examples and simple explanations. When you provide code examples, you always document your code thoroughly for maximum understanding"
	}
	hiYellowPrint("Ducky: ")
	yellowPrint("My role: '" + duckyRole + "'\n")

	conv = append(conv, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: duckyRole,
	})
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&model, "chat_model", "c", "gpt3", "Target model to chat with. Default is GPT3")
}
