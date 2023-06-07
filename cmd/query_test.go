package cmd

import (
	openai2 "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

func TestCanGetResponse(t *testing.T) {
	//prompt := "How can I quickly find a file and copy it to another location in bash shell?"
	prompt := "what command is best to change a filename in bash shell?"

	viper.SetEnvPrefix("ducky")
	viper.AutomaticEnv()
	apiKey := viper.GetString("DUCKY_API_KEY")
	resp, err := SendToGPT(apiKey, openai2.GPT3Dot5Turbo, []openai2.ChatCompletionMessage{{
		Role:    openai2.ChatMessageRoleUser,
		Content: prompt,
	}})

	if err != nil {
		t.Errorf("Error getting response from API: %v", err)
	}

	expected := "mv"
	if !strings.Contains(resp, expected) {
		t.Errorf("Expected %s to be in response %s", expected, resp)
	}
}

func TestFormattingOutputWithCode(t *testing.T) {
	//resp := "\n\nYou can use the 'find' command to search for the file and the 'cp' command to copy it to another location. Here is an example command:\n\n```bash\nfind /path/to/search -name \"filename.txt\" -exec cp {} /path/to/destination \\;\n```\n\nThis command will search for a file named 'filename.txt' in the directory '/path/to/search' and any subdirectories. Once the file is found, it will be copied to the directory '/path/to/destination'.\n\nYou can also use wildcards (*) in the file name to search for files that match a pattern. For example, to copy all files with a '.txt' extension from the current directory to a directory called 'textfiles', you can use the following command:\n\n```bash\nfind . -name \"*.txt\" -exec cp {} textfiles \\;\n```\n\nThis will find all files with a '.txt' extension in the current directory and any subdirectories, and copy them to the 'textfiles' directory."
	//resp := "\n\nThe command \"mv\" (short for \"move\") is best to change a filename in the bash shell. \n\nThe syntax for renaming a file is as follows: \n\n```\nmv filename1 filename2\n```\n\nHere, \"filename1\" is the name of the file you want to change, and \"filename2\" is the new name you want to give."
	resp := "Unfortunately, Go does not have a built-in do-while loop like some other programming languages. However, it is possible to simulate a do-while loop in Go using a for loop and a break statement.\n\nHere's an example of how to simulate a do-while loop in Go:\n\n```go\nfor {\n    // Code to be executed at least once\n\n    if condition {\n        // Code to be executed if condition is true\n\n        // Break out of the loop if necessary\n        break\n    }\n}\n```\n"
	codeSnippets := ExtractCodeSnippets(resp)

	err := PrintFormattedCodeSnippets(codeSnippets)
	if err != nil {
		t.Errorf("Error printing formatted code snippets: %v", err)
	}
}
