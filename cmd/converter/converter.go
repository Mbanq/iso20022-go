package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"

	pacs "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
)

func main() {
	// Read JSON file
	jsonFile, err := ioutil.ReadFile("sample_files/pacs.008.001.08_scenario1.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Unmarshal JSON to Go struct
	var doc pacs.Document
	if err := json.Unmarshal(jsonFile, &doc); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Marshal Go struct to XML
	xmlBytes, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling XML: %v", err)
	}

	// Add namespace to Document tag via string replacement
	xmlString := strings.Replace(string(xmlBytes), `<Document>`, `<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08">`, 1)

	// Add XML header
	output := []byte(xml.Header + xmlString)

	// Write XML to file
	err = ioutil.WriteFile("sample_files/pacs.008.001.08_scenario1.xml", output, 0644)
	if err != nil {
		log.Fatalf("Error writing XML file: %v", err)
	}

	log.Println("Successfully converted JSON to XML: sample_files/pacs.008.001.08_scenario1.xml")
}
