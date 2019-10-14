package crafter

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type GadgetItem struct {
	UriEntry
	Payload
}

func OutputAgent(listeningIp string, httpPort int, isOnlinePayload, shouldProxyRequest bool, gadgetsList []GadgetItem) {
	agentTemplateWithArguments, err := FillAgentWithArguments(listeningIp, httpPort, isOnlinePayload, shouldProxyRequest)
	if err != nil {
		log.Println("Error while reading the default ./agent/agent_plate.ps1 file")
		os.Exit(1)
	}

	var (
		payloadListArray []string
		payload          string
	)
	for i := 0; i < len(gadgetsList); i++ {
		if isOnlinePayload {
			payload = ""
		} else {
			payload = gadgetsList[i].PayloadData
		}
		payloadListArray = append(payloadListArray, `[pscustomobject]@{ UniqueId = '`+gadgetsList[i].UniqueId+`'; PayloadContent = '`+payload+`'; UriProtocol = '`+gadgetsList[i].UriProtocol+`'}`)
	}
	replacePayloadsPlaceholder := strings.NewReplacer("\"{{PAYLOADS}}\"", strings.Join(payloadListArray, ","))
	agentTemplateWithArguments = replacePayloadsPlaceholder.Replace(agentTemplateWithArguments)

	errW := ioutil.WriteFile("./output/agent.ps1", []byte(agentTemplateWithArguments), 0644)
	if errW != nil {
		log.Println("Error while writing the agent to file.")
	}
}

func FillAgentWithArguments(listeningIp string, httpPort int, isOnlinePayload, shouldProxyRequest bool) (string, error) {
	persistorTemplateInBytes, err := ioutil.ReadFile("./agent/agent_plate.ps1")
	persistorTemplateString := string(persistorTemplateInBytes)
	replacePlaceholders := strings.NewReplacer(
		"{listeningIp}", listeningIp,
		"{httpPort}", strconv.Itoa(httpPort),
		"\"{ISONLINEFETCH}\"", "$"+(strconv.FormatBool(isOnlinePayload)),
		"\"{PROXYREQUEST}\"", "$"+(strconv.FormatBool(shouldProxyRequest)))

	return replacePlaceholders.Replace(persistorTemplateString), err
}
