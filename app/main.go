package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
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

// What's left:
// write the JSON to a file via 'write-to-temp-then-rename into /run/wg-status.json' instead of just printing to stdout
// set up the cron job to run it every 15-30s.
// add wg interface back in?
func main() {
	// cmd := exec.Command("sudo", "wg", "show", "wg0", "dump")
	// stdout, err := cmd.Output()
	// if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// }

	line1 := "SOoqyUSZE72DBdG316aYsTIu/nChYAxTV6YLU8bzOVY=\tR6ioUf9Qxdr8BGwE+dM9g8dQdI8ynX/Ose0WRiie+jk=\t51820\toff\n"
	line2 := "+loDGuu6zAxLLwuCvdjJGy26l/q0UHI00co++uvDwXg=\t(none)\t192.24.141.184:51820\t10.6.0.2/32\t1788820051\t66418732\t867925452\toff\n"
	line3 := "yfjy63AUEw15fg76IabU4hHHSJO8qsmpuMdDC5VY73k=\t(none)\t(none)\t10.6.0.3/32\t0\t0\t0\toff\n"
	stdout := line1 + line2 + line3

	for _, deviceLine := range cleanedOutputTmp(stdout) {
		mappedDevice := mapDevice(deviceLine)
		devicesArray = append(devicesArray, mappedDevice)
	}

	jsonDevices, _ := json.MarshalIndent(devicesArray, "", " ")
	createTempFile(jsonDevices)
	fmt.Println(string(jsonDevices))
}

func createTempFile(stdout []byte) {
	file, err := os.CreateTemp("run", "wg-status.json.tmp")
	if err != nil {
		panic(err)
	}
	// defer os.Rename("run/wg-status.json.tmp*", "run/wg-status.json")

	if _, err := file.Write(stdout); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
	if err := os.Rename("run/wg-status.json.tmp*", "run/wg-status.json"); err != nil {
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

func cleanedOutputTmp(stdout string) []string {
	splitStr := s.Split(s.TrimSpace(stdout), "\n")
	return slices.Delete(splitStr, 0, 1)
}
