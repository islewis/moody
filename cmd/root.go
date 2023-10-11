/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"io/ioutil"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main.go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// set config info
	configName := "config"
	configType := "yaml"
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType(configType) // REQUIRED if the config file does not have the extension in the name

	// Get config path
	home, _ := os.UserHomeDir()
	dirPath := home + "/.config/moody"
	viper.AddConfigPath(dirPath)

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("Config not found, creating one...")

		// make sure dirs are created
		_ = os.MkdirAll(dirPath, os.ModePerm)

		// write file with any defaults
		defaultConfig := []byte(
			`checks:
  mood:
    description: "What was your mood today?"
    type: "range"
    min: "-2"
    max: "2"
  exercise:
    description: "Did you exercise today?"
    type: "bool"
  water:
    description: "How many glasses of water did you drink?"
    type: "int"`)

		fileToWrite := dirPath + "/" + configName + "." + configType
		err := ioutil.WriteFile(fileToWrite, defaultConfig, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	} else {
		return
	}
}
