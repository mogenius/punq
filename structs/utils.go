package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mogenius/punq/logger"

	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
)

const PingSeconds = 10

func MarshalUnmarshal(datagram *Datagram, data interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := json.Marshal(datagram.Payload)
	if err != nil {
		datagram.Err = err.Error()
		return
	}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		datagram.Err = err.Error()
	}
}

func PrettyPrint(i interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	iJson, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		logger.Log.Error(err.Error())
	}
	PrettyPrintJSON(iJson)
}

func PrettyPrintString(i interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	iJson, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return string(iJson)
}

func PrettyPrintJSON(input []byte) {
	var raw json.RawMessage
	err := json.Unmarshal(input, &raw)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, raw, "", "  ")
	if err != nil {
		logger.Log.Error(err.Error())
	}

	colorizeString(buf.String())
}

func colorizeString(prettyData string) {
	var rawData interface{}
	err := json.Unmarshal([]byte(prettyData), &rawData)
	if err != nil {
		log.Fatalln("Error unmarshalling JSON:", err)
	}

	colorize(rawData, "", false)
	fmt.Println()
}

func colorize(data interface{}, prefix string, isKey bool) {
	switch v := data.(type) {
	case map[string]interface{}:
		colorizeMap(v, prefix)
	case []interface{}:
		colorizeArray(v, prefix)
	case string:
		if isKey {
			color.FgLightCyan.Printf("%s\"%s\"", prefix, v)
		} else {
			color.FgLightYellow.Printf("\"%s\"", v)
		}
	default:
		fmt.Print(v)
	}
}

func colorizeMap(m map[string]interface{}, prefix string) {
	fmt.Print("{\n")
	newPrefix := prefix + "  "
	for k, v := range m {
		colorize(k, newPrefix, true)
		color.FgWhite.Print(": ")
		colorize(v, newPrefix, false)
		color.FgWhite.Print(",\n")
	}
	fmt.Printf("%s}", prefix)
}

func colorizeArray(a []interface{}, prefix string) {
	fmt.Print("[\n")
	newPrefix := prefix + "  "
	for _, v := range a {
		colorize(v, newPrefix, false)
		color.FgWhite.Print(",\n")
	}
	fmt.Printf("%s]", prefix)
}

func MilliSecSince(since time.Time) int64 {
	return time.Since(since).Milliseconds()
}

func MicroSecSince(since time.Time) int64 {
	return time.Since(since).Microseconds()
}

func DurationStrSince(since time.Time) string {
	duration := MilliSecSince(since)
	durationStr := fmt.Sprintf("%d ms", duration)
	if duration <= 0 {
		duration = MicroSecSince(since)
		durationStr = fmt.Sprintf("%d μs", duration)
	}
	return durationStr
}
