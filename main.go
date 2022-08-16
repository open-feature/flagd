package main

import (
	"github.com/open-feature/flagd/cmd"
	"flag"
	"fmt"
	"os"
)
var (
    version = "dev"
    commit  = "HEAD"
    date    = "unknown"
)

func main() {

	showVersion:= flag.Bool("version", false, "show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("flagd %s (%s) built at %s\n", version, commit, date)
		os.Exit(0)
	} else {
		cmd.Execute()
	}
}
