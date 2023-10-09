package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"strconv"
	"time"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Record mood, and other daily checks",
	Long:  "TO DO",
	Run: func(cmd *cobra.Command, args []string) {
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

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(checkinCmd)
}

func initConfig() {
	// set config info
	configName := "config"
	configType:= ".yaml"
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType(configType)   // REQUIRED if the config file does not have the extension in the name

	// Get config path
	home, _ := os.UserHomeDir()
	dirPath := home + "/.config/moody"
	viper.AddConfigPath(dirPath)

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// make sure dirs are created
		fmt.Println("Config not found!")
		fmt.Println(viper.ConfigFileUsed)
		_ = os.MkdirAll(dirPath, os.ModePerm)
		// write file with any defaults
		// Init config file here when ready to build out

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

		fileToWrite := dirPath + "/" + configName + configType
		err := ioutil.WriteFile(fileToWrite, defaultConfig, 0644)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
	} else {
		return
		// Config file was found but another error was produced
	}
}
