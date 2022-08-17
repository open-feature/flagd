package main

import (
	"github.com/open-feature/flagd/cmd"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
