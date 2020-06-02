package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"github.com/stetsd/monk-conf"
	"os"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup app",
	Long:  `setup app`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.EnvParseToConfigMap()

		if err != nil {
			panic(err)
		}

		location, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		rows := []string{
			"flyway.url=jdbc:postgresql://" + config.Get("DB_HOST") + ":" + config.Get("DB_PORT") + "/" + config.Get("DB_NAME"),
			"flyway.user=" + config.Get("DB_USER"),
			"flyway.password=" + config.Get("DB_PASS"),
			"flyway.locations=filesystem:" + location + "/migrations",
		}

		file, err := os.Create("flyway.conf")

		if err != nil {
			panic(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()

		writer := bufio.NewWriter(file)

		for _, row := range rows {
			_, err := writer.WriteString(row + "\n")
			if err != nil {
				panic(err)
			}
		}

		if err := writer.Flush(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
