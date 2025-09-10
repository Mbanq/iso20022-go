package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Mbanq/iso20022-go/pkg/fednow"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file.xml>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Read the XML file
	xmlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		os.Exit(1)
	}

	// Call the parser with the file content
	jsonOutput, err := fednow.Parse(xmlFile)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput))
}
