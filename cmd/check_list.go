package cmd

import (
	"fmt"
	"log"
	"database/sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var checkListCmd = &cobra.Command{
	Use: "list",
	Short: "List all daily checks",
	Run: func(cmd *cobra.Command, args []string) {
		// Find and read the config file
		err := viper.ReadInConfig()
		if err != nil {panic(err)}

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

		fmt.Println("List of checks:")
		for key := range checksList {
			fmt.Println(key)
			fmt.Sprintf("\n %s: %s", checksList[key].Name, checksList[key].Manifest["description"])
			
			todaysDate := time.Now()
			formattedDate := todaysDate.Format("20060102")
  		dbSelect := fmt.Sprintf(`SELECT %s FROM DailyEntries WHERE "date" = %s`, key, formattedDate)

			db, err := sql.Open("sqlite3", "~/.config/moody/data.db")
			if err != nil {log.Fatal(err)}
			dbReturn, err := db.Exec(dbSelect)
			defer db.Close()

			if err != nil {
				fmt.Println("Todays input: ")
			} else {
				fmt.Sprintf("Todays input: %s", dbReturn)
			}
		}
	},
}

func init() {
	checkCmd.AddCommand(checkListCmd)
}
