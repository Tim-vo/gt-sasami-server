package cmd

import (
	"fmt"

	conf "github.com/Tim-vo/gt-sasami-server/config"
	cli "github.com/spf13/cobra"
)

// Version command
func init() {
	rootCmd.AddCommand(&cli.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println(conf.Executable + " - " + conf.GitVersion)
		},
	})
}
