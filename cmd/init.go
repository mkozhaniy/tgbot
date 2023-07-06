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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create database",
	Long: `Create tables(users, orders ...), for creating need enviroment variable
	TG_BOT_DB, it is URL for database, also you need TG_BPT_LOGPATH, for logging`,
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
		db.InitTables(dbConn)
		if err != nil {
			logger.Warn(err.Error())
			panic(err)
		}
		sqlDb, err := dbConn.DB()
		if err != nil {
			logger.Warn(err.Error())
			panic(err)
		}
		sqlDb.Close()
		logger.Sync()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
