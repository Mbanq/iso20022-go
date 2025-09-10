package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Mbanq/iso20022-go/pkg/fednow/bah"
	"github.com/Mbanq/iso20022-go/pkg/fednow/config"
	"github.com/Mbanq/iso20022-go/pkg/fednow/pacs"
)

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("failed to load config: %s\n", err)
		return
	}

	jsonFile, err := ioutil.ReadFile("schema_Cct_ex.json")
	if err != nil {
		fmt.Printf("failed to read json file: %s\n", err)
		return
	}

	// Unmarshal the json into a FedNowMessage struct
	var fednowMessage pacs.FedNowMessageCCT
	if err := json.Unmarshal(jsonFile, &fednowMessage); err != nil {
		fmt.Printf("Error unmarshalling json: %s\n", err)
		return
	}

	pacsDoc, err := pacs.BuildPacs008(jsonFile, config)
	if err != nil {
		fmt.Printf("Error building pacs.008: %s\n", err)
		return
	}

	// Generate BAH
	bahDoc, err := bah.BuildBah(string(fednowMessage.FedNowMsg.Identifier.MessageID), config, "pacs.008.001.08")
	if err != nil {
		fmt.Printf("Error building bah: %s\n", err)
		return
	}

	// Marshal pacs.008 to XML
	pacsOutput, err := xml.MarshalIndent(pacsDoc, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling pacs.008 to xml: %v\n", err)
		return
	}
	finalPacsOutput := strings.Replace(string(pacsOutput), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08\">", 1)

	err = ioutil.WriteFile("output.xml", []byte(finalPacsOutput), 0644)
	if err != nil {
		fmt.Println("Error writing pacs.008 XML file:", err)
		return
	}
	fmt.Println("Successfully generated output.xml")

	// Marshal BAH to XML
	bahOutput, err := xml.MarshalIndent(bahDoc, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling bah to xml: %v\n", err)
		return
	}
	finalBahOutput := strings.Replace(string(bahOutput), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	finalBahOutput = strings.Replace(finalBahOutput, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	err = ioutil.WriteFile("bah_output.xml", []byte(finalBahOutput), 0644)
	if err != nil {
		fmt.Println("Error writing BAH XML file:", err)
		return
	}
	fmt.Println("Successfully generated bah_output.xml")
}
