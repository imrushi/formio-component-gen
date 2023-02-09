package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	componentLocation := flag.String("location", "component", "Create Component at location")
	jsonFileName := flag.String("json", "data.json", "Load json from which component should be generated")

	flag.Parse()

	// Open the JSON file and parse the story in it.
	f, err := os.Open(*jsonFileName)
	if err != nil {
		panic(err)
	}
	compo, err := JsonDecode(f)
	if err != nil {
		panic(err)
	}

	err = Generate(compo.Data, *componentLocation)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
