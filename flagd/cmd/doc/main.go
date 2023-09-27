package main

import (
	"log"

	"github.com/open-feature/flagd/flagd/cmd"
)

const docPath = "../docs/reference/flagd-cli"

func main() {
	if err := cmd.GenerateDoc(docPath); err != nil {
		log.Fatal(err)
	}
}
