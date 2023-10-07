package cmd

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Record mood, and other daily checks",
	Long: "TO DO",
	Run: func(cmd *cobra.Command, args []string) {
		moodValues := []int{-2,-1,0,1,2}
		fmt.Println("On a scale of (-2,-1, 0, 1, 2) how was your mood today?")
		fmt.Print("Mood: ")
		var dailyMood int
		fmt.Scanln(&dailyMood)
		if slices.Contains(moodValues, dailyMood) != true {
			fmt.Println("")
			fmt.Printf("Invalid mood value. You printed %s, value needs to be one of -2, -1, 0, 1, 2\n", strconv.Itoa(dailyMood))
		} else {
			fmt.Println("")
			fmt.Printf("Marked today as a %s\n", strconv.Itoa(dailyMood))
		}
	},
}

func init() {
	rootCmd.AddCommand(checkinCmd)
}
