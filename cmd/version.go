package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of FlagD",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		details, ok := debug.ReadBuildInfo()
		if ok && details.Main.Version != "(devel)" {
			Version = details.Main.Version
			for _, i := range details.Settings {
				if i.Key == "vcs.time" {
					Date = i.Value
				}
				if i.Key == "vcs.revision" {
					Commit = i.Value
				}
			}
		}
		fmt.Printf("flagd %s (%s) built at %s\n", Version, Commit, Date)
	},
}
