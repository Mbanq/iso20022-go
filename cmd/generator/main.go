package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Mbanq/iso20022-go/pkg/fednow"
	"github.com/Mbanq/iso20022-go/pkg/fednow/config"
	"github.com/Mbanq/iso20022-go/pkg/fednow/pacs"
)

func main() {
	xsdPath := flag.String("xsd", "", "Path to the XSD file")
	messageId := flag.String("messageId", "", "Message ID (e.g., pacs.008.001.08)")
	inputPath := flag.String("in", "", "Path to the input JSON file")
	outPath := flag.String("out", "", "Path to the output XML file")
	flag.Parse()

	if *xsdPath == "" || *messageId == "" || *inputPath == "" || *outPath == "" {
		fmt.Println("Usage: go run main.go --xsd <path> --messageId <id> --in <path> --out <path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Generating XML for message ID '%s' using XSD '%s'\n", *messageId, *xsdPath)

	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("failed to load config: %s\n", err)
		return
	}

	jsonFile, err := ioutil.ReadFile(*inputPath)
	if err != nil {
		fmt.Printf("failed to read json file: %s\n", err)
		return
	}

	var fednowMessage fednow.FedNowMessage

	switch *messageId {
	case "pacs.008.001.08":
		var msg pacs.FedNowMessageCCT
		if err := json.Unmarshal(jsonFile, &msg); err != nil {
			fmt.Printf("Error unmarshalling json for pacs.008: %s\n", err)
			return
		}
		fednowMessage = msg
	case "pacs.002.001.10":
		var msg pacs.FedNowMessageACK
		if err := json.Unmarshal(jsonFile, &msg); err != nil {
			fmt.Printf("Error unmarshalling json for pacs.002: %s\n", err)
			return
		}
		fednowMessage = msg
	case "pacs.004.001.10":
		var msg pacs.FedNowMessageRtn
		if err := json.Unmarshal(jsonFile, &msg); err != nil {
			fmt.Printf("Error unmarshalling json for pacs.004: %s\n", err)
			return
		}
		fednowMessage = msg
	default:
		fmt.Printf("unsupported message type: %s\n", *messageId)
		return
	}

	xmlData, err := fednow.Generate(*xsdPath, *messageId, config, fednowMessage)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(*outPath, xmlData, 0644)
	if err != nil {
		fmt.Printf("Error writing XML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated XML file at '%s'\n", *outPath)
}
