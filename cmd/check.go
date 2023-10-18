package cmd

import (
	"fmt"
	"os"
  "strings"
	"database/sql"
	"time"
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

		// handle "yes" and "no" inputs as bools
		if strings.ToLower(inputStr) == "yes" {
			inputStr = "true"
		} else if strings.ToLower(inputStr) == "no" {
			inputStr = "false"
		}

		_, err := strconv.ParseBool(inputStr)
		if err != nil {
			log.Fatal(err)
		}
		return inputStr
	}
	// Error out before returning the placeholder string
	log.Fatal(`Error: type "` + r.Manifest["type"].(string) + `" for check ` + r.Name + "not recognized")
	return ""
}

func (r Check) WriteToDB(value, pathToDB string) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {log.Fatal(err)}

	// Add an empty date entry which we'll ALTER as input comes in
	todaysDate := time.Now()
	formattedDate := todaysDate.Format("20060102")
  dateInsert := fmt.Sprintf(`INSERT INTO DailyEntries ("date") VALUES ("%s");`, formattedDate)
	_, _ = db.Exec(dateInsert)
	
	// Add Column. If it exists, we get an error which we dont do anything with 
	// Could handle this error more gracefully
	dbTable := fmt.Sprintf(`ALTER TABLE DailyEntries ADD COLUMN %s`, r.Name)
	_, _ = db.Exec(dbTable)

  dbAlter := fmt.Sprintf(`
		UPDATE DailyEntries
		SET "%s" = "%s"
		WHERE date = "%s";`, r.Name, value, formattedDate)
	_, err = db.Exec(dbAlter)
	if err != nil {log.Fatal(err)}

	defer db.Close()
}

var checkCmd = &cobra.Command{
	Use: "check",
	Short: "Fill out your daily checks",
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
			value := checksList[key].TakeInput()
			// need to pull the DB location from some sort of default config
			home, _ := os.UserHomeDir()
			checksList[key].WriteToDB(value, home + "/.config/moody/data.db")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
