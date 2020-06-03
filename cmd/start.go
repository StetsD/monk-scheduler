package cmd

import (
	"github.com/spf13/cobra"
	config "github.com/stetsd/monk-conf"
	"github.com/stetsd/monk-scheduler/internal/app"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start http server",
	Long:  `start http server`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.EnvParseToConfigMap()

		if err != nil {
			panic(err)
		}

		server := app.NewApp(conf)
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
