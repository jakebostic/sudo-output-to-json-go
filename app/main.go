package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"slices"
	"strconv"
	s "strings"
	"time"

	h "github.com/dustin/go-humanize"
)

type Device struct {
	DeviceName      string `json:"device_name"`
	AllowedIps      string `json:"allowed_ips"`
	Endpoint        string `json:"endpoint"`
	LatestHandshake string `json:"latest_handshake"`
	ListenPort      string `json:"listen_port"`
	MbReceived      string `json:"mb_received"`
	MbTransferred   string `json:"mb_transferred"`
}

var devicesArray []Device
var knownDeviceMap = map[string]string{
	"+loDGuu6zAxLLwuCvdjJGy26l/q0UHI00co++uvDwXg=": "iPhone",
	"yfjy63AUEw15fg76IabU4hHHSJO8qsmpuMdDC5VY73k=": "Desktop"}

func main() {
	cmd := exec.Command("sudo", "wg", "show", "wg0", "dump")
	stdout, err := cmd.Output()
	if err != nil {
			fmt.Println(err.Error())
			return
	}
	for _, deviceLine := range cleanedOutput(stdout) {
		mappedDevice := mapDevice(deviceLine)
		devicesArray = append(devicesArray, mappedDevice)
	}
	jsonDevices, _ := json.MarshalIndent(devicesArray, "", " ")
	fmt.Println(string(jsonDevices))
}

func mapDevice(deviceLine string) Device {
	formattedLine := s.Fields(deviceLine) // returns slice with line fields as substrings
	return Device{
		deviceName(formattedLine[0]),
		formattedLine[3],
		formattedLine[2],
		parseTime(formattedLine[4]),
		"51820",
		parseTransferBytes(formattedLine[5]),
		parseTransferBytes(formattedLine[6])}
}

func deviceName(publicKey string) string {
	if knownDeviceName, publicKeyMatches := knownDeviceMap[publicKey]; publicKeyMatches {
		return knownDeviceName
	}
	return "Unknown"
}

func parseTime(latestHandshake string) string {
	handshake := parseStrToInt(latestHandshake)
	return h.Time(time.Unix(handshake, 0).UTC())
}

func parseTransferBytes(mbData string) string {
	transferInt := parseStrToInt(mbData)
	return h.BigIBytes(big.NewInt(transferInt))
}

func parseStrToInt(value string) int64 {
	parsedInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return parsedInt
}

func cleanedOutput(stdout []byte) []string {
	splitStr := s.Split(s.TrimSpace(string(stdout)), "\n")
	return slices.Delete(splitStr, 0, 1)
}
