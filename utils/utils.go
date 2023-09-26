package utils

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	jsoniter "github.com/json-iterator/go"
)

const APP_NAME = "k8s"

func Pointer[K any](val K) *K {
	return &val
}

type ResponseError struct {
	Error string `json:"error,omitempty"`
}

type Release struct {
	TagName   string `json:"tag_name"`
	Published string `json:"published_at"`
}

func IsProduction() bool {
	stage := os.Getenv("stage")
	if stage == "" {
		stage = os.Getenv("STAGE")
	}
	return Equals([]string{"prod", "production"}, strings.ToLower(stage))
}

func Equals(s []string, str string) bool {
	for _, v := range s {
		if str == v {
			return true
		}
	}
	return false
}

func IsNewReleaseAvailable() bool {
	latestRelease := "https://api.github.com/repos/mogenius/punq/releases/latest"
	resp, err := http.Get(latestRelease)
	if err != nil {
		logger.Log.Errorf("Error getting release: %s", err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Log.Errorf("failed to fetch latest release: %s", string(bodyBytes))
		return false
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		logger.Log.Errorf("Error decoding release: %s", err.Error())
		return false
	}

	fmt.Printf("Your version:      %s\n", version.Ver)
	fmt.Printf("Available version: %s        (published: %s ago)\n", release.TagName, JsonStringToHumanDuration(release.Published))

	if strings.Contains(release.TagName, version.Ver) {
		fmt.Println("You are up-to-date 🥰.")
		return false
	} else {
		fmt.Println("Your version is outdated 😭!\n❗️Please update punq: https://punq.dev\n")
		return true
	}
}

func CreateError(err error) ResponseError {
	return ResponseError{
		Error: err.Error(),
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func Diff(a []string, b []string) []string {
	diff := make([]string, 0)

	if len(a) != len(b) {
		return a
	}

	// Create a map to store the count of each string in array 'a'
	countMap := make(map[string]int)
	for _, str := range a {
		countMap[str]++
	}

	// Check if all strings in array 'b' are present in the map
	for _, str := range b {
		count, ok := countMap[str]
		if !ok || count == 0 {
			diff = append(diff, str)
		} else {
			countMap[str]--
		}
	}

	// Add any remaining items in countMap to the diff slice
	for str, count := range countMap {
		if count > 0 {
			diff = append(diff, str)
		}
	}

	return diff
}

func ContainsInt(v int, a []int) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println(fmt.Errorf("error while opening browser, %v", err))
	}
}

func ConfirmTask(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}

func HashString(data string) string {
	var buf bytes.Buffer

	// Serialize the object
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256(buf.Bytes())

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}

func SelectIndexInteractive(question string, noOfElements int) int {
	for {
		// Prompt the user for an index
		fmt.Printf("\nSelect a number between (1-%d) (or type 'exit' or 'all'): ", noOfElements)
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			return -1
		}
		if input == "all" {
			return -2
		}

		// Try to convert the user input into an integer index
		index, err := strconv.Atoi(input)
		if err != nil || index <= 0 || index > noOfElements {
			fmt.Println("Invalid index. Please try again.")
			continue
		}
		return index
	}
}

func FillWith(s string, targetLength int, chars string) string {
	if len(s) >= targetLength {
		return TruncateText(s, targetLength)
	}
	for i := 0; len(s) < targetLength; i++ {
		s = s + chars
	}

	return s
}

func TruncateText(s string, max int) string {
	if max < 4 || max > len(s) {
		return s
	}
	return s[:max-4] + " ..."
}

func FunctionName() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	return runtime.FuncForPC(counter).Name()
}

func ParseJsonStringArray(input string) []string {
	val := []string{}
	var jsonOnSteroids = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := jsonOnSteroids.Unmarshal([]byte(input), &val); err != nil {
		logger.Log.Errorf("jsonStringArrayToStringArray: Failed to parse: '%s' to []string.", input)
	}
	return val
}

func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func HttpHeader(additionalName string) http.Header {
	return http.Header{
		"x-app":         []string{fmt.Sprintf("%s%s", APP_NAME, additionalName)},
		"x-app-version": []string{version.Ver}}
}

func CreateDirIfNotExist(dir string) {
	_, err := os.Stat(dir)

	// If directory does not exist create it
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0755)
		if errDir != nil {
			logger.Log.Error(err.Error())
		}
	}
}

func DeleteDirIfExist(dir string) {
	_, err := os.Stat(dir)

	// If directory does not exist create it
	if os.IsExist(err) {
		errDir := os.RemoveAll(dir)
		if errDir != nil {
			logger.Log.Error(err.Error())
		}
	}
}

// parseIPs parses a slice of IP address strings into a slice of net.IP.
func parseIPs(ips []string) ([]net.IP, error) {
	var parsed []net.IP
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return nil, fmt.Errorf("invalid IP address: %s", ip)
		}
		parsed = append(parsed, parsedIP.To4())
	}
	return parsed, nil
}
