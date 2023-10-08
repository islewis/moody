package cmd

import (
	"fmt"
	"time"
	"log"
	"strconv"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Record mood, and other daily checks",
	Long: "TO DO",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Daily Mood
		moodValues := []int{-2,-1,0,1,2}
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
				INSERT INTO DailyEntries (date, mood) VALUES (`+ formattedDate +`, `+ strconv.Itoa(dailyMood) +`);
			`
			_, err = db.Exec(dbq)
    			if err != nil {
        			log.Fatal(err)
    			}

			defer db.Close()
			var mood int
			var date string
			err = db.QueryRow("SELECT date, mood FROM DailyEntries WHERE date = 20231008").Scan(&date, &mood)
			fmt.Println("Todays date: "+ date+ ". Todays mood: "+ strconv.Itoa(mood))

			fmt.Println("")
			fmt.Printf("Marked today as a %s\n", strconv.Itoa(dailyMood))
		}
	},
}

func init() {
	rootCmd.AddCommand(checkinCmd)
}
