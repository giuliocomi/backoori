package crafter

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Payloads struct {
	Payloads []Payload `json:"payloads"`
}

type Payload struct {
	PayloadName string `json:"payload_name"`
	PayloadData string `json:"payload_template"`
	UniqueId    string //for make agent know which URL to fetch
}

func GetAvailablePayloads(jsonPayloads []byte) Payloads {
	var payloads = Payloads{}

	errParse := json.Unmarshal(jsonPayloads, &payloads)
	if errParse != nil {
		log.Fatal("Error parsing the JSON file containing the payloads to use")
	}
	return payloads
}

func (payload *Payload) GetPayloadParams() []string {
	regex := regexp.MustCompile("{{\\w*}}") // find the params to set by searching for the placeholder  {{}}
	matches := regex.FindAllString(payload.PayloadData, -1)

	return matches
}

func (payload *Payload) SetPayloadParams(paramsToFill, paramsFilled []string) {
	var replacePlaceholders *strings.Replacer
	for index, _ := range paramsToFill {
		replacePlaceholders = strings.NewReplacer(paramsToFill[index], paramsFilled[index])
		payload.PayloadData = replacePlaceholders.Replace(payload.PayloadData)
		if replacePlaceholders != nil {
		} else {
			log.Println("Error replacing place holders with user-specified parameters")
		}
	}
}

func DeployCradleGadgetPayload(templatePayload Payload) string {
	//set a unique random value to identify the gadget payload
	templatePayload.UniqueId = GenerateUniquePayloadId()
	//create web resource with payloadID as filename and payload data as content
	payloadContent := []byte(templatePayload.PayloadData)
	errCreateGadgetFile := ioutil.WriteFile("./output/gadgets/"+templatePayload.UniqueId, payloadContent, 0600)
	if errCreateGadgetFile != nil {
		log.Println(errCreateGadgetFile)
	}
	return templatePayload.UniqueId
}

func GenerateUniquePayloadId() string {
	var randomID = make([]byte, 24)
	_, _ = rand.Read(randomID)
	return fmt.Sprintf("%x", randomID)
}
