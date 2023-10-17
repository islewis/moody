package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Record mood, as well as and other daily checks",
	Long:  "TO DO",
	Run: func(cmd *cobra.Command, args []string) {
		// Loop through checks
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {
			panic(err)
		}

		// Get list of checks
		checkMap := viper.GetStringMap("checks")
		var checks []string
		for key := range checkMap {
			checks = append(checks, key)
		}

		// Iterate through checks
		for _, check := range checks {
			// Validate input differently for each data type
			if viper.GetString("checks."+check+".type") == "range" {
				handleRangeCheck(check)
			}
			if viper.GetString("checks."+check+".type") == "int" {
				handleRangeCheck(check)
			}
		}

		// Get Daily Mood
		moodValues := []int{-2, -1, 0, 1, 2}
		fmt.Println("On a scale of (-2,-1, 0, 1, 2) how was your mood today?")
		fmt.Print("Mood: ")
		var dailyMood int
		fmt.Scanln(&dailyMood)
		// Get date
		todaysDate := time.Now()
		formattedDate := todaysDate.Format("20060102")
		// Verify its an acceptable value
		if slices.Contains(moodValues, dailyMood) != true {
			fmt.Println("")
			fmt.Printf("Invalid mood value. You printed %s, value needs to be one of -2, -1, 0, 1, 2\n", strconv.Itoa(dailyMood))
		} else {
			db, err := sql.Open("sqlite3", ":memory:")
			if err != nil {
				log.Fatal(err)
			}
			dbq := `
				CREATE TABLE DailyEntries (date PRIMARY KEY, mood INTEGER);
				INSERT INTO DailyEntries (date, mood) VALUES (` + formattedDate + `, ` + strconv.Itoa(dailyMood) + `);
			`
			_, err = db.Exec(dbq)
			if err != nil {
				log.Fatal(err)
			}

			defer db.Close()
			var mood int
			var date string
			err = db.QueryRow("SELECT date, mood FROM DailyEntries WHERE date = 20231008").Scan(&date, &mood)
			fmt.Println("Todays date: " + date + ". Todays mood: " + strconv.Itoa(mood))

			fmt.Println("")
			fmt.Printf("Marked today as a %s\n", strconv.Itoa(dailyMood))
		}
	},
}

func handleRangeCheck(name string) {

	// Before we take an input, confirm we have a maximum and a minimum defined in the config
	maxValLoc := ("checks." + name + ".max")
	if viper.GetString(maxValLoc) == "" {
		errorMessage := "ERROR: the " + name + " definition is missing a maximum value.\n"
		errorMessage += "All 'range' types need to have a max and a min defined. This should be " + maxValLoc + " in your config\n."
		log.Fatalf(errorMessage)
	}

	minValLoc := ("checks." + name + ".min")
	if viper.GetString(maxValLoc) == "" {
		errorMessage := "ERROR: the " + name + " definition is missing a minimum value.\n"
		errorMessage += "All 'range' types need to have a max and a min defined. This should be " + minValLoc + " in your config\n."
		log.Fatalf(errorMessage)
	}

	// Get input until we have a value that is within our valid range
	for {
		fmt.Println(viper.GetString("checks." + name + ".description"))
		fmt.Print(name + ": ")
		var input int
		fmt.Scanln(&input)
		//input, err := strconv.Atoi(input)
		//if err != nil {
		//	log.Fatal(err)
		//}

		maxVal, _ := strconv.Atoi(viper.GetString(maxValLoc))
		minVal, _ := strconv.Atoi(viper.GetString(minValLoc))
		if input > maxVal {
			fmt.Println("\n")
			fmt.Println("Input is greater than the maximum of " + viper.GetString(maxValLoc) + ", try again.")
			fmt.Println("\n")
			continue
		} else if input < minVal {
			fmt.Println("Input is less than than the minimum of " + viper.GetString(minValLoc) + ", try again.")
			continue
		} else {
			break
		}
	}

}

func init() {
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
