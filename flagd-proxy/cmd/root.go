package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	Version string
	Commit  string
	Date    string
	Debug   bool
)

var rootCmd = &cobra.Command{
	Use:               "flagd",
	Short:             "flagd-proxy allows flagd to subscribe to CRD changes without the required permissions.",
	Long:              ``,
	DisableAutoGenTag: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString(logFormatFlagName) == "console" {
			banner.InitString(colorable.NewColorableStdout(), true, true, `
	{{ .AnsiColor.BrightRed }}	 ______   __       ________   _______    ______      
	{{ .AnsiColor.BrightRed }}	/_____/\ /_/\     /_______/\ /______/\  /_____/\     
	{{ .AnsiColor.BrightRed }}	\::::_\/_\:\ \    \::: _  \ \\::::__\/__\:::_ \ \    
	{{ .AnsiColor.BrightRed }}	 \:\/___/\\:\ \    \::(_)  \ \\:\ /____/\\:\ \ \ \   
	{{ .AnsiColor.BrightRed }}	  \:::._\/ \:\ \____\:: __  \ \\:\\_  _\/ \:\ \ \ \  
	{{ .AnsiColor.BrightRed }}	   \:\ \    \:\/___/\\:.\ \  \ \\:\_\ \ \  \:\/.:| | 
	{{ .AnsiColor.BrightRed }}	    \_\/     \_____\/ \__\/\__\/ \_____\/   \____/_/ 
	{{ .AnsiColor.BrightRed }}	                                   Kubernetes Proxy  
{{ .AnsiColor.Default }}
`)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, commit string, date string) {
	Version = version
	Commit = commit
	Date = date
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "x", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agent.yaml)")
	rootCmd.AddCommand(startCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".agent" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".agent")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "using config file:", viper.ConfigFileUsed())
	}
}
