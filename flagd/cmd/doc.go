package cmd

import (
	"fmt"

	"github.com/spf13/cobra/doc"
)

// GenerateDoc generates cobra docs of the cmd
func GenerateDoc(path string) error {
	linkHandler := func(name string) string {
		return name
	}

	filePrepender := func(filename string) string {
		return "<!-- markdownlint-disable-file -->\n<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->\n"
	}

	if err := doc.GenMarkdownTreeCustom(rootCmd, path, filePrepender, linkHandler); err != nil {
		return fmt.Errorf("error generating docs: %w", err)
	}
	return nil
}
