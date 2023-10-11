package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type Check struct {
	Name     string
	Manifest map[string]interface{}
}

func (r Check) Type() string {
	manifest := r.Manifest
	return manifest["type"].(string)
}

func (r Check) TakeInput() string {
	// Handle Input seperately for each data type
	if r.Type() == "range" {
		// Before we take an range input, confirm we have a maximum and a minimum defined in the config
		if r.Manifest["max"].(string) == "" {
			errorMessage := "ERROR: the " + r.Name + " definition is missing a maximum value.\n"
			errorMessage += "All 'range' types need to have a max and a min defined. This should be checks." + r.Name + ".max in your config\n."
			log.Fatalf(errorMessage)
		}

		if r.Manifest["min"].(string) == "" {
			errorMessage := "ERROR: the " + r.Name + " definition is missing a minimum value.\n"
			errorMessage += "All 'range' types need to have a max and a min defined. This should be checks." + r.Name + ".min in your config\n."
			log.Fatalf(errorMessage)
		}

		// Get input until we have a value that is within our valid range
		for {
			fmt.Print("\n")
			fmt.Println(r.Manifest["description"])
			fmt.Print(r.Name + ": ")
			var inputStr string
			var inputInt int
			fmt.Scanln(&inputStr)

			// Confirm our input is really an int
			inputInt, err := strconv.Atoi(inputStr)
			if err != nil {
				log.Fatal(err)
			}

			maxVal, _ := strconv.Atoi(r.Manifest["max"].(string))
			minVal, _ := strconv.Atoi(r.Manifest["min"].(string))
			if inputInt > maxVal {
				fmt.Println("\n")
				fmt.Println("Input is greater than the maximum of " + strconv.Itoa(maxVal) + ", try again.")
				continue
			} else if inputInt < minVal {
				fmt.Println("\n")
				fmt.Println("Input is less than than the minimum of " + strconv.Itoa(minVal) + ", try again.")
				continue
			} else {
				return inputStr
			}
		}
	} else if r.Type() == "int" {
		fmt.Print("\n")
		fmt.Println(r.Manifest["description"])
		fmt.Print(r.Name + ": ")
		var inputStr string
		fmt.Scanln(&inputStr)

		// Confirm our input is really an int
		_, err := strconv.Atoi(inputStr)
		if err != nil {
			log.Fatal(err)
		}

		return inputStr
	} else if r.Type() == "bool" {
		fmt.Print("\n")
		fmt.Println(r.Manifest["description"])
		fmt.Print(r.Name + ": ")
		var inputStr string
		fmt.Scanln(&inputStr)

		_, err := strconv.ParseBool(inputStr)
		if err != nil {
			log.Fatal(err)
		}
		return inputStr
	} else {
		log.Fatal(`Error: type "` + r.Manifest["type"].(string) + `" for check ` + r.Name)
	}
	return ""
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		// Loop through checks
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {
			panic(err)
		}

		// Get list of checks
		checksConfig := viper.GetStringMap("checks")
		checksList := make([]Check, 0)
		for key := range checksConfig {
			check := Check{
				Name:     key,
				Manifest: checksConfig[key].(map[string]interface{}),
			}
			checksList = append(checksList, check)
		}
		for key := range checksList {
			checksList[key].TakeInput()
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
