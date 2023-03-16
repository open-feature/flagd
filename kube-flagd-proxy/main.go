package main

import (
	"bytes"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
	"github.com/open-feature/flagd/kube-flagd-proxy/cmd"
)

var (
	version    = "dev"
	commit     = "HEAD"
	date       = "unknown"
	bannerText = `
	{{ .AnsiColor.BrightRed }}	 ______   __       ________   _______    ______      
	{{ .AnsiColor.BrightRed }}	/_____/\ /_/\     /_______/\ /______/\  /_____/\     
	{{ .AnsiColor.BrightRed }}	\::::_\/_\:\ \    \::: _  \ \\::::__\/__\:::_ \ \    
	{{ .AnsiColor.BrightRed }}	 \:\/___/\\:\ \    \::(_)  \ \\:\ /____/\\:\  \ \ \   
	{{ .AnsiColor.BrightRed }}	  \:::._\/ \:\ \____\:: __  \ \\:\\_  _\/ \:\ \ \ \  
	{{ .AnsiColor.BrightRed }}	   \:\ \    \:\/___/\\:.\ \  \ \\:\_\ \ \  \:\/.:| | 
	{{ .AnsiColor.BrightRed }}	    \_\/     \_____\/ \__\/\__\/ \_____\/   \____/_/ 
	{{ .AnsiColor.BrightRed }}	                                   Kubernetes Proxy  
{{ .AnsiColor.Default }}
`
)

func main() {
	banner.Init(colorable.NewColorableStdout(), true, true,
		bytes.NewBufferString(bannerText))
	cmd.Execute(version, commit, date)
}
