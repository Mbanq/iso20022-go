// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Mbanq/iso20022-go/pkg/fednow"
	"github.com/Mbanq/iso20022-go/pkg/fednow/config"
	"github.com/Mbanq/iso20022-go/pkg/fednow/pacs"
)

func main() {
	// Load configuration
	configFile, err := ioutil.ReadFile("../../config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	var msgConfig config.Config
	if err := json.Unmarshal(configFile, &msgConfig); err != nil {
		log.Fatalf("Failed to unmarshal config json: %v", err)
	}

	// Read the input JSON file
	jsonFile, err := ioutil.ReadFile("../../schema_Rtn_ex.json")
	if err != nil {
		log.Fatalf("Failed to read json file: %v", err)
	}

	// Create a new pacs.004 message from the JSON data
	var pacs004Msg pacs.FedNowMessageRtn
	if err := json.Unmarshal(jsonFile, &pacs004Msg); err != nil {
		log.Fatalf("Failed to unmarshal pacs.004 message from json: %v", err)
	}

	// Generate the pacs.004.001.10 message
	xsdPath := "../../Internal/XSD/iso/pacs.004.001.10.xsd"
	xmlBytes, err := fednow.Generate(xsdPath, "pacs.004.001.10", &msgConfig, pacs004Msg)
	if err != nil {
		log.Fatalf("Failed to generate pacs.004 message: %v", err)
	}

	// Define the output file path
	outputDir := "sample_files"
	outputFile := filepath.Join(outputDir, "Output_pacs004.xml")

	// Create the sample_files directory if it doesn't exist
	if err = os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write the XML to the output file
	if err = ioutil.WriteFile(outputFile, xmlBytes, 0644); err != nil {
		log.Fatalf("Failed to write xml file: %v", err)
	}

	fmt.Println("Successfully generated Output_pacs004.xml")
}
