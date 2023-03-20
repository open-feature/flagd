package cmd

import (
	"strings"

	"github.com/spf13/cobra/doc"
)

// GenerateDoc generates cobra docs of the cmd
func GenerateDoc(path string) error {
	linkHandler := func(name string) string {
		return strings.ReplaceAll(name, ".md", "")
	}

	filePrepender := func(filename string) string {
		return "<!-- markdownlint-disable-file -->\n"
	}
	return doc.GenMarkdownTreeCustom(rootCmd, path, filePrepender, linkHandler)
}
