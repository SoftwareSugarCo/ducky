/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	model        string
	defaultModel = openai.GPT3Dot5Turbo
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [prompt]",
	Short: "Query the OpenAI API with a prompt",
	Run:   handleQueryCmd,
}

func handleQueryCmd(cmd *cobra.Command, args []string) {
	prompt := args[0]
	apiKey := viper.GetString("api_key")

	if apiKey == "" {
		fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
		return
	}

	// Default model selected is GPT3.5 Turbo
	model = GetModel(model)
	fmt.Println("Using model: " + model)

	// Send out the query to OpenAI
	resp, err := SendToGPT(apiKey, model, []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract any code snippets from the response and store them
	codeSnippets := ExtractCodeSnippets(resp)

	if ToFile {
		fmt.Println("ToFile mode not yet implemented")
	}

	fmt.Println(resp)
	fmt.Println(TITLEDIVSTR)
	_ = PrintFormattedCodeSnippets(codeSnippets)
	fmt.Println(DIVSTR)
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&model, "query_model", "q", "gpt3", "Target model to send query to. Default is GPT3.")
}
