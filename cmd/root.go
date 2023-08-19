package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	ToFile     bool
	StreamMode bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ducky",
	Short: "Ducky is a CLI tool to interact with OpenAI's ChatGPT",
	Long: `Ducky is a CLI tool to interact with OpenAI's ChatGPT. It is optimized for programming questions.
	Ducky will extract the code snippet and format it in the interpreted language. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ducky.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&ToFile, "to_file", "f", false, "Save each snippet to a file with the language inferred from the snippet.")
	rootCmd.PersistentFlags().BoolVarP(&StreamMode, "stream", "s", false, "Stream the output instead of returning it all at once... this is how chatGPT responds.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ducky" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ducky")
	}
	viper.SetEnvPrefix("ducky")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Config file not found, check for environment variables
		if !viper.IsSet("api_key") {
			fmt.Println("No API key found. Please set the DUCKY_API_KEY environment variable or add it to your config file.")
			os.Exit(1)
		}
	}
}
