package utils

import (
	"bufio"
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"punq/logger"
	"punq/version"
	"runtime"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const APP_NAME = "k8s"

var YamlTemplatesFolder embed.FS

func Pointer[K any](val K) *K {
	return &val
}

type ResponseError struct {
	Error string `json:"error,omitempty"`
}

type Volume struct {
	Namespace  NamespaceDisplayName `json:"namespace"`
	VolumeName string               `json:"volumeName"`
	SizeInGb   int                  `json:"sizeInGb"`
}

type NamespaceDisplayName struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func MountPath(namespaceName string, volumeName string, defaultReturnValue string) string {
	if CONFIG.Kubernetes.RunInCluster {
		return fmt.Sprintf("%s/%s_%s", CONFIG.Misc.DefaultMountPath, namespaceName, volumeName)
	} else {
		pwd, err := os.Getwd()
		pwd += "/temp"
		if err != nil {
			logger.Log.Errorf("StatsMogeniusNfsVolume PWD Err: %s", err.Error())
		} else {
			return pwd
		}
	}
	return defaultReturnValue
}

func StorageClassForClusterProvider(clusterProvider string) string {
	var nfsStorageClassStr string = "default"
	// TODO: "DOCKER_ENTERPRISE", "DOKS", "LINODE", "IBM", "ACK", "OKE", "OPEN_SHIFT"
	switch clusterProvider {
	case "EKS":
		nfsStorageClassStr = "gp2"
	case "GKE":
		nfsStorageClassStr = "standard-rwo"
	case "AKS":
		nfsStorageClassStr = "default"
	case "OTC":
		nfsStorageClassStr = "csi-disk"
	case "BRING_YOUR_OWN":
		nfsStorageClassStr = "default"
	default:
		logger.Log.Errorf("CLUSTERPROVIDER '%s' HAS NOT BEEN TESTED YET! Returning 'default'.", clusterProvider)
		nfsStorageClassStr = "default"
	}
	return nfsStorageClassStr
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
		fmt.Errorf("error while opening browser, %v", err)
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
		"x-authorization":  []string{CONFIG.Kubernetes.ApiKey},
		"x-cluster-mfa-id": []string{CONFIG.Kubernetes.ClusterMfaId},
		"x-app":            []string{fmt.Sprintf("%s%s", APP_NAME, additionalName)},
		"x-app-version":    []string{version.Ver},
		"x-cluster-name":   []string{CONFIG.Kubernetes.ClusterName}}
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

func GetVolumeMountsForK8sManager() ([]Volume, error) {
	result := []Volume{}

	// Create an http client
	client := &http.Client{}

	// Create a new request using http
	url := fmt.Sprintf("%s/storage/k8s/cluster-project-storage/list", CONFIG.ApiServer.Http_Server)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	// Add headers to the http request
	req.Header = HttpHeader("")
	// TODO: REMOVE - THIS IS JUST FOR DEBUGGING
	// if CONFIG.Misc.Debug && CONFIG.Misc.Stage == "local" {
	// 	req.Header["x-authorization"] = []string{"mo_7bf5c2b5-d7bc-4f0e-b8fc-b29d09108928_0hkga6vjum3p1mvezith"}
	// 	req.Header["x-cluster-mfa-id"] = []string{"a141bd85-c986-402c-9475-5bdc4679293b"}
	// }

	// Send the request and get a response
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, &result)
	return result, err
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

// // FindSmallestSubnet finds the smallest subnet that includes all given IP addresses.
// func FindSmallestSubnet(ipStrings []string) *net.IPNet {
// 	ips, err := parseIPs(ipStrings)
// 	if err != nil {
// 		fmt.Println("Error parsing IP addresses:", err)
// 		return nil
// 	}

// 	sort.Slice(ips, func(i, j int) bool {
// 		return bytes.Compare(ips[i], ips[j]) < 0
// 	})
// 	minIP, maxIP := ips[0], ips[len(ips)-1]

// 	mask := net.CIDRMask(commonPrefixLen(minIP, maxIP), 32)
// 	return &net.IPNet{IP: minIP, Mask: mask}
// }

// func LastIpMinusOne(network *net.IPNet) net.IP {
// 	var ip net.IP
// 	for i := 0; i < len(network.IP); i++ {
// 		ip = append(ip, network.IP[i]|(^network.Mask[i]))
// 	}
// 	if ip4 := ip.To4(); ip4 != nil {
// 		ip4[3]--
// 		return ip4
// 	}
// 	ip[15]--
// 	return ip
// }

// // commonPrefixLen finds the length of the common prefix of a and b in bits.
// func commonPrefixLen(a, b net.IP) (cpl int) {
// 	for i := 0; i < 4; i++ {
// 		diff := uint(a[i] ^ b[i])
// 		for diff != 0 {
// 			diff >>= 1
// 			cpl++
// 		}
// 	}
// 	return 32 - cpl
// }
