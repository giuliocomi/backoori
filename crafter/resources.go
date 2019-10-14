package crafter

import (
	"io/ioutil"
	"log"
)

func LoadResources(jsonURIprotocols, jsonPayloads string) (UriList, Payloads) {
	uriList, errOpen := ioutil.ReadFile(jsonURIprotocols)
	if errOpen != nil {
		log.Fatal("Error reading the JSON file containing URI protocols")
	}
	payloads, errOpen := ioutil.ReadFile(jsonPayloads)
	if errOpen != nil {
		log.Fatal("Error reading the JSON file containing the payloads to use")
	}
	uriProtocols := GetURIs(uriList)
	availablePayloads := GetAvailablePayloads(payloads)
	return uriProtocols, availablePayloads
}
