package crafter

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	stdin = bufio.NewReader(os.Stdin)
)

func HelpMenu(version string) {
	fmt.Println("Backoori" + version + ": tool aided persistence via Windows URI schemes abuse")
	fmt.Println("Generate a ready-to-launch Powershell agent that will backdoor specific Universal URI Apps with fileless payloads of your choice.")
}

func WebServerDialog() (string, int, int) {
	var listeningAddress string
	var port, timeout int

	for {
		fmt.Println("Provide the IP of the machine where the web server to deliver the payloads is hosted:")
		_, errA := fmt.Scanf("%s", &listeningAddress)
		FlushInputStream(stdin)
		if errA != nil || net.ParseIP(listeningAddress) == nil {
			log.Println("Incorrect IPv4 address")
			continue
		}
		fmt.Println("Provide the port of the machine were the web server to deliver the payloads is hosted:")
		_, errP := fmt.Scanf("%d", &port)
		FlushInputStream(stdin)
		if errP != nil || port < 0 || port > 65535 {
			log.Println("Incorrect port value")
			continue
		}
		fmt.Println("Set timeout (in seconds) before closing the connection:")
		_, errT := fmt.Scanf("%d", &timeout)
		FlushInputStream(stdin)
		if errT != nil || timeout < 0 {
			log.Println("Incorrect timeout amount")
			continue
		}
		fmt.Printf("Deploying for %d seconds the HTTP server to deliver the gadgets on %s:%d\n", timeout, listeningAddress, port)
		return listeningAddress, port, timeout
	}
}

func PayloadDialog(payloadsToDisplay Payloads) int {
	var payloadIndex int

	for {
		//list payloads
		fmt.Println("Payloads loaded via JSON file:")
		for index := 0; index < len(payloadsToDisplay.Payloads); index++ {
			fmt.Printf("%d) %s\n", index, payloadsToDisplay.Payloads[index].PayloadName)
		}
		//user chooses index (starting from 0) of the payload
		fmt.Print("Enter the index of the payload to use: ")
		_, err := fmt.Scanln(&payloadIndex)
		if err != nil || uint(len(payloadsToDisplay.Payloads)) <= uint(payloadIndex) {
			FlushInputStream(stdin)
			log.Println("Selected payload index not available")
			continue
		}
		fmt.Println("Done, new gadget ready")
		return payloadIndex
	}
}

func ParamsDialog(payload Payload) Payload {
	var (
		paramsToFill, paramsFilled []string
		paramToDisplay             string
	)
	paramsToFill = payload.GetPayloadParams()
	paramsFilled = make([]string, len(paramsToFill))

	for index, _ := range paramsToFill {
		paramToDisplay = strings.Trim(paramsToFill[index], "{{")
		paramToDisplay = strings.Trim(paramToDisplay, "}}")
		fmt.Printf("Specify value for parameter %s\n", paramToDisplay)
		_, errP := fmt.Scanln(&paramsFilled[index])
		if errP != nil {
			FlushInputStream(stdin)
		}
	}
	// replace payload with every parameter filled
	payload.SetPayloadParams(paramsToFill, paramsFilled)
	return payload
}

func FlagDialog(isOnlinePayload, shouldProxyRequest string) (bool, bool) {
	isOnlineBool, err1 := strconv.ParseBool(isOnlinePayload)
	shouldProxyBool, err2 := strconv.ParseBool(shouldProxyRequest)
	if err1 != nil || err2 != nil {
		log.Println("Cannot convert flags passed via argument to Boolean")
		os.Exit(1)
	}
	return isOnlineBool, shouldProxyBool
}

func UriToBackdoorDialog(uriList UriList) int {
	var uriIndex int

	for {
		//list payloads
		fmt.Println("URI protocols loaded via JSON file:")
		for index := 0; index < len(uriList.UriList); index++ {
			fmt.Printf("%d) %s\n", index, uriList.UriList[index].UriProtocol)
		}
		//user chooses number of the payload
		fmt.Print("Enter the index of the URI protocol to backdoor: ")
		_, err := fmt.Scanln(&uriIndex)
		if err != nil || uint(len(uriList.UriList)) <= uint(uriIndex) {
			FlushInputStream(stdin)
			fmt.Println("Selected URI protocol not available in the JSON config file. Add it as an entry first.")
			continue
		}
		return uriIndex
	}
}

func ChooseAnotherUriToBackdoorDialog() bool {
	fmt.Println("Press 'c' for preparing another gadget, any key otherwise to exit")
	var key string
	fmt.Scanln(&key)
	if string(key) != string('c') {
		fmt.Println("'c' not pressed, exiting menu")
	}
	return string(key) == string('c')
}

func OnExitDialog(isOnlinePayload bool, ip string, port, timeout int) {
	fmt.Println("Payloads and Agent forged and ready. Agent has been written to ./output/agent.ps1.")
	if isOnlinePayload {
		fmt.Printf("The webserver was started at: %s:%d \n", ip, port)
		ch := make(chan bool, 1)
		defer close(ch)

		exitSignal := make(chan os.Signal)
		signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)
		timer := time.NewTimer(time.Duration(timeout) * time.Second)
		defer timer.Stop()
		select {
		case <-exitSignal:
			log.Println("Backoori terminated.")
		case <-timer.C:
			fmt.Println("Timeout for web server connection has been reached. Quitting, Bye.")
			cleanedFolder := OnWebServerShutdown()
			if cleanedFolder {
				fmt.Println("gadgets folder successfully cleaned.")
			} else {
				log.Println("Failed clearing gadgets folder upon exit.")
			}
		}
	} else {
		fmt.Println("The payloads have been directly embedded in the agent for offline use. Quitting, Bye.")
		os.Exit(0)
	}
}

func FlushInputStream(r *bufio.Reader, a ...interface{}) {
	_, _ = r.Discard(r.Buffered())
	_, _ = fmt.Fscanln(r, a...)
}