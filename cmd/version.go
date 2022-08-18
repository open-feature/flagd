package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of FlagD",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("flagd %s (%s) built at %s\n", Version, Commit, Date)
	},
}
