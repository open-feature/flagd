package cmd

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of FlagD",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "dev" {
			details, ok := debug.ReadBuildInfo()
			if ok && details.Main.Version != "" && details.Main.Version != "(devel)" {
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
		}
		log.Printf("flagd %s (%s) built at %s\n", Version, Commit, Date)
	},
}
