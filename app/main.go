package main

import (
	"encoding/json"
	"math/big"
	"os"
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

var knownDeviceMap = map[string]string{
	"phone_key_****": "iPhone",
	"desktop_key_****": "Desktop"}

const statusDir = "/run" // change to "run" for laptop/desktop testing

func main() {
	cmd := exec.Command("sudo", "wg", "show", "wg0", "dump")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var devicesArray []Device
	for _, deviceLine := range cleanedOutput(stdout) {
		mappedDevice := mapDevice(deviceLine)
		devicesArray = append(devicesArray, mappedDevice)
	}
	jsonDevices, _ := json.MarshalIndent(devicesArray, "", " ")
	buildJsonFile(jsonDevices)
}

func buildJsonFile(stdout []byte) {
	file, err := os.CreateTemp(statusDir, "config.json.tmp")
	if err != nil {
		panic(err)
	}

	if err := file.Chmod(0644); err != nil {
		panic(err)
	}

	if _, err := file.Write(stdout); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	newStatusDir := statusDir + "/config.json"
	if err := os.Rename(file.Name(), newStatusDir); err != nil {
		panic(err)
	}
}

func mapDevice(deviceLine string) Device {
	formattedLine := s.Fields(deviceLine)
	return Device{
		deviceName(formattedLine[0]),
		formattedLine[3],
		formattedLine[2],
		parseTime(formattedLine[4]),
		"****",
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

func cleanedOutputTmp(stdout string) []string {
	splitStr := s.Split(s.TrimSpace(stdout), "\n")
	return slices.Delete(splitStr, 0, 1)
}
