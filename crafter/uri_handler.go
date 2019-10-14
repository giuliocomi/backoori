package crafter

import (
	"encoding/json"
	"log"
)

type UriList struct {
	UriList []UriEntry `json:"URI_Protocols"`
}

type UriEntry struct {
	UriProtocol string `json:"URI_Protocol"`
}

// returns the struct with info for AppX IDs and URI schemes to provide to the backoori-agent
func GetURIs(configFile []byte) UriList {

	var uri = UriList{}
	errParse := json.Unmarshal(configFile, &uri)
	if errParse != nil {
		log.Fatal("Error parsing the JSON file containing the list of URI protocols")
	}
	return uri
}
