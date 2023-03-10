package main

import (
	"log"

	"github.com/open-feature/flagd/flagd/cmd"
)

const docPath = "../docs/configuration"

func main() {
	if err := cmd.GenerateDoc(docPath); err != nil {
		log.Fatal(err)
	}
}
