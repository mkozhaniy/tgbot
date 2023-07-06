/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tgbot/internal/db"
	"go.uber.org/zap"
)

// droptbCmd represents the droptb command
var droptbCmd = &cobra.Command{
	Use:   "droptb",
	Short: "drop all tables",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		url := os.Getenv("TG_BOT_DB")
		logPaths := os.Getenv("TG_BOT_LOGPATH")
		dbConn, err := db.Start(url)
		cfg := zap.NewDevelopmentConfig()
		cfg.OutputPaths = []string{logPaths}
		logger := zap.Must(cfg.Build())
		if err != nil {
			logger.Warn(err.Error())
			panic(err)
		}
		defer logger.Sync()
		db.DropTables(dbConn)
		if err != nil {
			logger.Warn(err.Error())
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(droptbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// droptbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// droptbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
