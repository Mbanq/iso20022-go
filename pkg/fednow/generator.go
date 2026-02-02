package fednow

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	admi007 "github.com/mbanq/iso20022-go/ISO20022/admi_007_001_01"
	camt029 "github.com/mbanq/iso20022-go/ISO20022/camt_029_001_09"
	camt056 "github.com/mbanq/iso20022-go/ISO20022/camt_056_001_08"
	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	pacs002 "github.com/mbanq/iso20022-go/ISO20022/pacs_002_001_10"
	pacs004 "github.com/mbanq/iso20022-go/ISO20022/pacs_004_001_10"
	pacs008 "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
	pain013 "github.com/mbanq/iso20022-go/ISO20022/pain_013_001_07"
	"github.com/mbanq/iso20022-go/pkg/common"
	"github.com/mbanq/iso20022-go/pkg/fednow/admi"
	bah "github.com/mbanq/iso20022-go/pkg/fednow/bah"
	"github.com/mbanq/iso20022-go/pkg/fednow/camt"
	config "github.com/mbanq/iso20022-go/pkg/fednow/config"
	"github.com/mbanq/iso20022-go/pkg/fednow/pacs"
	"github.com/mbanq/iso20022-go/pkg/fednow/pain"
)

type xsdCacheEntry struct {
	rootElement    string
	messageElement string
	wrapperElement string
	rootNs         string
}

var (
	xsdCache = make(map[string]*xsdCacheEntry)
	cacheMux = &sync.RWMutex{}
)

// Generate creates a FedNow XML envelope for a given message ID using the specified XSD file.
func Generate(xsdPath, messageType string, config *config.Config, message FedNowMessage) ([]byte, error) {
	// First, try to read from the cache with a read lock.
	cacheMux.RLock()
	entry, found := xsdCache[messageType]
	cacheMux.RUnlock()

	if !found {

		cacheMux.Lock()
		defer cacheMux.Unlock()
		entry, found = xsdCache[messageType]
		if !found {
			rootElement, messageElement, wrapperElement, rootNs, err := findWrapperForMessageID(xsdPath, messageType)
			if err != nil {
				return nil, fmt.Errorf("error finding wrapper element: %v", err)
			}

			entry = &xsdCacheEntry{
				rootElement:    rootElement,
				messageElement: messageElement,
				wrapperElement: wrapperElement,
				rootNs:         rootNs,
			}
			xsdCache[messageType] = entry
		}
	}

	handler, ok := messageHandlers[messageType]
	if !ok {
		return nil, fmt.Errorf("unsupported message type: %s", messageType)
	}

	bah, pacs, err := handler(config, message)
	if err != nil {
		return nil, err
	}

	// Construct the final XML
	finalXML := fmt.Sprintf("<%s xmlns=\"%s\">\n    <%s>\n        <%s>\n%s\n%s\n        </%s>\n    </%s>\n</%s>",
		entry.rootElement,
		entry.rootNs,
		entry.messageElement,
		entry.wrapperElement,
		bah,
		pacs,
		entry.wrapperElement,
		entry.messageElement,
		entry.rootElement,
	)

	return []byte(finalXML), nil
}

type messageHandler func(cfg *config.Config, msg FedNowMessage) (string, string, error)

var messageHandlers = map[string]messageHandler{
	"admi.007.001.01": handleAdmi007,
	"pacs.008.001.08": handlePacs008,
	"pacs.002.001.10": handlePacs002,
	"pacs.004.001.10": handlePacs004,
	"pain.013.001.07": handlePain013,
	"camt.056.001.08": handleCamt056,
	"camt.029.001.09": handleCamt029,
}

func handleAdmi007(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(admi.FedNowMessageRctAck)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for admi.007.001.01")
	}

	appHdr, document, err := GenerateAdmi007("admi.007.001.01", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	admi007Doc := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:admi.007.001.01\">", 1)

	return bah, admi007Doc, nil
}

func handlePacs002(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(pacs.FedNowMessageACK)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for pacs.002.001.10")
	}

	appHdr, document, err := GeneratePacs002("pacs.002.001.10", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	pacs002 := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10\">", 1)

	return bah, pacs002, nil
}

func handleCamt056(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(camt.FedNowMessageCxlReq)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for camt.056.001.08")
	}

	appHdr, document, err := GenerateCamt056("camt.056.001.08", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	camt056Doc := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:camt.056.001.08\">", 1)

	return bah, camt056Doc, nil
}

func handleCamt029(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(camt.FedNowMessageCxlRsp)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for camt.029.001.09")
	}

	appHdr, document, err := GenerateCamt029("camt.029.001.09", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	camt029Doc := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:camt.029.001.09\">", 1)

	return bah, camt029Doc, nil
}

func handlePacs004(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(pacs.FedNowMessageRtn)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for pacs.004.001.10")
	}

	appHdr, document, err := GeneratePacs004("pacs.004.001.10", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	pacs004 := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:pacs.004.001.10\">", 1)

	return bah, pacs004, nil
}

func handlePain013(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(pain.FedNowMessageRFP)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for pain.013.001.07")
	}

	appHdr, document, err := GeneratePain013("pain.013.001.07", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	pain013 := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:pain.013.001.07\">", 1)

	return bah, pain013, nil
}

func GeneratePacs004(messageType string, msgConfig *config.Config, message pacs.FedNowMessageRtn) (*head.BusinessApplicationHeaderV02, *pacs004.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := pacs.BuildPacs004Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func GeneratePacs002(messageType string, msgConfig *config.Config, message pacs.FedNowMessageACK) (*head.BusinessApplicationHeaderV02, *pacs002.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := pacs.BuildPacs002Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func GenerateCamt056(messageType string, msgConfig *config.Config, message camt.FedNowMessageCxlReq) (*head.BusinessApplicationHeaderV02, *camt056.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := camt.BuildCamt056Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func GenerateCamt029(messageType string, msgConfig *config.Config, message camt.FedNowMessageCxlRsp) (*head.BusinessApplicationHeaderV02, *camt029.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := camt.BuildCamt029Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func GenerateAdmi007(messageType string, msgConfig *config.Config, message admi.FedNowMessageRctAck) (*head.BusinessApplicationHeaderV02, *admi007.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := admi.BuildAdmi007Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func handlePacs008(cfg *config.Config, message FedNowMessage) (string, string, error) {
	msg, ok := message.(pacs.FedNowMessageCCT)
	if !ok {
		return "", "", fmt.Errorf("invalid message type for pacs.008.001.08")
	}

	appHdr, document, err := GeneratePacs008("pacs.008.001.08", cfg, msg)
	if err != nil {
		return "", "", err
	}

	appHdrPayload, err := xml.MarshalIndent(appHdr, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling AppHdr: %v", err)
	}

	bah := strings.Replace(string(appHdrPayload), "<BusinessApplicationHeaderV02>", "<AppHdr xmlns=\"urn:iso:std:iso:20022:tech:xsd:head.001.001.02\">", 1)
	bah = strings.Replace(bah, "</BusinessApplicationHeaderV02>", "</AppHdr>", 1)

	documentPayload, err := xml.MarshalIndent(document, "            ", "    ")
	if err != nil {
		return "", "", fmt.Errorf("error marshalling document: %v", err)
	}

	pacs008 := strings.Replace(string(documentPayload), "<Document>", "<Document xmlns=\"urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08\">", 1)

	return bah, pacs008, nil
}

func GeneratePacs008(messageType string, msgConfig *config.Config, message pacs.FedNowMessageCCT) (*head.BusinessApplicationHeaderV02, *pacs008.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := pacs.BuildPacs008Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

func GeneratePain013(messageType string, msgConfig *config.Config, message pain.FedNowMessageRFP) (*head.BusinessApplicationHeaderV02, *pain013.Document, error) {

	now := time.Now().In(common.EstLocation)
	// Override creation date and time with current EST time
	message.FedNowMsg.CreationDateTime = common.ISODateTime(now)

	appHdr, err := bah.BuildBah(string(message.FedNowMsg.Identifier.MessageID), msgConfig, messageType)
	if err != nil {
		return nil, nil, err
	}

	document, err := pain.BuildPain013Struct(message, msgConfig)
	if err != nil {
		return nil, nil, err
	}

	return appHdr, document, nil
}

// findWrapperForMessageID dynamically parses the XSD to find the correct wrapper element.
func findWrapperForMessageID(xsdPath string, messageId string) (string, string, string, string, error) {
	xsdFile, err := os.Open(xsdPath)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error opening XSD file: %v", err)
	}
	defer xsdFile.Close()

	decoder := xml.NewDecoder(xsdFile)
	var rootElement, messageElement, wrapperElementName, rootNs string
	nsMap := make(map[string]string)
	var elementStack []xml.StartElement
	var choiceRefs []string

	// First pass: find the root, message, and choice refs
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", "", "", fmt.Errorf("error decoding XSD token (pass 1): %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			elementStack = append(elementStack, se)
			if se.Name.Local == "schema" {
				for _, attr := range se.Attr {
					if attr.Name.Space == "xmlns" {
						nsMap[attr.Name.Local] = attr.Value
					} else if attr.Name.Local == "targetNamespace" {
						rootNs = attr.Value
					}
				}
			} else if se.Name.Local == "element" {
				parent, grandParent := getParents(elementStack)
				if parent.Name.Local == "schema" {
					if rootElement == "" {
						rootElement = getAttrName(se, "name")
					}
				} else if parent.Name.Local == "sequence" && grandParent.Name.Local == "complexType" {
					ggParent := getGrandparent(elementStack)
					if ggParent.Name.Local == "element" && getAttrName(ggParent, "name") == rootElement {
						minOccurs := getAttrName(se, "minOccurs")
						if minOccurs != "0" {
							ref := getAttrName(se, "ref")
							if messageElement == "" && ref != "" {
								messageElement = ref
							}
						}
					}
				} else if parent.Name.Local == "choice" {
					ref := getAttrName(se, "ref")
					if ref != "" {
						choiceRefs = append(choiceRefs, ref)
					}
				}
			}
		case xml.EndElement:
			if len(elementStack) > 0 {
				elementStack = elementStack[:len(elementStack)-1]
			}
		}
	}

	// Second pass: find the wrapper element that references the correct document namespace
	xsdFile.Seek(0, 0)
	decoder = xml.NewDecoder(xsdFile)
	var currentWrapper string
	var inWrapper bool

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", "", "", fmt.Errorf("error decoding XSD token (pass 2): %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "element" && isChoiceRef(getAttrName(se, "name"), choiceRefs) {
				currentWrapper = getAttrName(se, "name")
				inWrapper = true
			} else if inWrapper && se.Name.Local == "element" {
				ref := getAttrName(se, "ref")
				if strings.Contains(ref, ":") {
					parts := strings.Split(ref, ":")
					refPrefix := parts[0]
					if ns, ok := nsMap[refPrefix]; ok {
						if strings.Contains(ns, messageId) {
							wrapperElementName = currentWrapper
							goto end_pass_2 // Found it, exit loop
						}
					}
				}
			}
		case xml.EndElement:
			if inWrapper && se.Name.Local == currentWrapper {
				inWrapper = false
				currentWrapper = ""
			}
		}
	}
end_pass_2:

	if rootElement == "" || messageElement == "" || wrapperElementName == "" {
		return "", "", "", "", fmt.Errorf("could not determine all required elements (root: '%s', message: '%s', wrapper: '%s')", rootElement, messageElement, wrapperElementName)
	}

	return rootElement, messageElement, wrapperElementName, rootNs, nil
}

func getParents(stack []xml.StartElement) (xml.StartElement, xml.StartElement) {
	if len(stack) < 2 {
		return xml.StartElement{}, xml.StartElement{}
	}
	if len(stack) < 3 {
		return stack[len(stack)-2], xml.StartElement{}
	}
	return stack[len(stack)-2], stack[len(stack)-3]
}

func getGrandparent(stack []xml.StartElement) xml.StartElement {
	if len(stack) < 4 {
		return xml.StartElement{}
	}
	return stack[len(stack)-4]
}

func getAttrName(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func isChoiceRef(name string, refs []string) bool {
	for _, ref := range refs {
		if ref == name {
			return true
		}
	}
	return false
}
