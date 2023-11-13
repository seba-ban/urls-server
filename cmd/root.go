/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/seba-ban/urls/server"
	"github.com/spf13/cobra"
)

var doMigrate bool
var migrationsFolder string
var sqliteDbPath string

var rootCmd = &cobra.Command{
	Use:   "urls-server",
	Short: "Server for managing saved urls",
	Run: func(cmd *cobra.Command, args []string) {
		if doMigrate {
			log.Println("Running migrations")
			server.Migrate(sqliteDbPath, migrationsFolder)
		}
		log.Println("Starting server")
		server.Run(sqliteDbPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&doMigrate, "migrate", "m", false, "run migrations")
	rootCmd.Flags().StringVarP(&migrationsFolder, "migrations-folder", "f", "", "folder with migrations")
	rootCmd.MarkFlagDirname("migrations-folder")
	if doMigrate {
		rootCmd.MarkFlagRequired("migrations-folder")
	}
	rootCmd.Flags().StringVarP(&sqliteDbPath, "sqlite-db-path", "d", "", "path to sqlite db")
	rootCmd.MarkFlagRequired("sqlite-db-path")
}
