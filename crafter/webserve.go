package crafter

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	listeningAddress string
	listeningPort    string
	schema           = bytes.Buffer{}
)

//TODO: https, payload encryption and pass AES decryption key to persistor agent that will decrypt the payloads before deploying the gadgets
func SetupWebServer(ip net.IP, port int) {
	listeningAddress = ip.String()
	listeningPort = strconv.Itoa(port)

	cleanGadgetsFolder()
	folder := http.FileServer(http.Dir("./output/gadgets"))
	http.Handle("/", folder)
	schema.WriteString(strings.Join([]string{listeningAddress, listeningPort}, ":"))
	errWebServer := http.ListenAndServe(schema.String(), nil)
	if errWebServer != nil {
		log.Println(errWebServer)
	}
}

func OnWebServerShutdown() bool {
	return cleanGadgetsFolder()
}

func cleanGadgetsFolder() bool {
	var (
		errC, errD error
		files      []string
		dir        string = "./output/gadgets/"
	)
	files, errD = filepath.Glob(filepath.Join(dir, "*"))
	for _, file := range files {
		errC = os.RemoveAll(file)
	}
	return errC == nil && errD == nil
}

// TODO: add connection logging and show which gadget id has been downloaded