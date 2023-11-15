package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

const (
	BOOL   = "boolean"
	STRING = "string"
)

/*
A simple random feature flag generator for testing purposes. Output is saved to "random.json".

Configurable options:

	-c : feature flag count (ex:go run ff_gen.go -c 500)
	-t : type of feature flag (ex:go run ff_gen.go -t string). Support "boolean" and "string"
*/
//nolint:gosec
func main() {
	// Get flag count
	var flagCount int
	flag.IntVar(&flagCount, "c", 100, "Number of flags to generate")

	// Get flag type : Boolean, String
	var flagType string
	flag.StringVar(&flagType, "t", BOOL, "Type of flags to generate")

	flag.Parse()

	if flagType != STRING && flagType != BOOL {
		fmt.Printf("Invalid type %s. Falling back to default %s", flagType, BOOL)
		flagType = BOOL
	}

	root := Flags{}
	root.Flags = make(map[string]Flag)

	switch flagType {
	case BOOL:
		root.setBoolFlags(flagCount)
	case STRING:
		root.setStringFlags(flagCount)
	}

	bytes, err := json.Marshal(root)
	if err != nil {
		fmt.Printf("Json error: %s ", err.Error())
		return
	}

	err = os.WriteFile("./random.json", bytes, 0o444)
	if err != nil {
		fmt.Printf("File write error: %s ", err.Error())
		return
	}
}

func (f *Flags) setBoolFlags(toGen int) {
	for i := 0; i < toGen; i++ {
		variant := make(map[string]any)
		variant["on"] = true
		variant["off"] = false

		f.Flags[fmt.Sprintf("flag%d", i)] = Flag{
			State:          "ENABLED",
			DefaultVariant: randomSelect("on", "off"),
			Variants:       variant,
		}
	}
}

func (f *Flags) setStringFlags(toGen int) {
	for i := 0; i < toGen; i++ {
		variant := make(map[string]any)
		variant["key1"] = "value1"
		variant["key2"] = "value2"

		f.Flags[fmt.Sprintf("flag%d", i)] = Flag{
			State:          "ENABLED",
			DefaultVariant: randomSelect("key1", "key2"),
			Variants:       variant,
		}
	}
}

type Flags struct {
	Flags map[string]Flag `json:"flags"`
}

type Flag struct {
	State          string         `json:"state"`
	DefaultVariant string         `json:"defaultVariant"`
	Variants       map[string]any `json:"variants"`
}

//nolint:gosec
func randomSelect(chooseFrom ...string) string {
	return chooseFrom[rand.Intn(len(chooseFrom))]
}
