/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tgbot/internal/db"
	"github.com/tgbot/tgbot"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tgbot",
	Short: "starting bot",
	Long: `Start listening updates from bot and sending 
	messages and response on callbacks, you need enviroment variables,
	TG_BOT_DB, TG_BOT_LOGPATH, TG_BOT_TOKEN, URL database, logs path and 
	bot token respectively`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		exit := make(chan struct{}, 1)
		go start_bot(os.Getenv("TG_BOT_DB"), os.Getenv("TG_BOT_LOGPATH"),
			os.Getenv("TG_BOT_TOKEN"), exit)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Write \"stop\" if you want stop the bot")
		for scanner.Scan() {
			if scanner.Text() == "stop" {
				exit <- struct{}{}
				break
			}
		}
	},
}

func Execute() {
	rootCmd.Execute()
}

func start_bot(url string, logpath string, token string, ch chan struct{}) {
	cfg := zap.NewDevelopmentConfig()
	logPaths := logpath
	cfg.OutputPaths = []string{logPaths}
	logger := zap.Must(cfg.Build())
	dbConn, err := db.Start(url)
	if err != nil {
		logger.Sugar().Panic(err.Error())
		panic(err)
	}
	source := db.Database{
		Db:     dbConn,
		Logger: logger,
	}
	tgbot.Start(logger, source, source, source, source, token, ch)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tgbot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
