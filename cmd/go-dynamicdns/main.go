package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const IP_LOG_PATH = "./previous-ips.log"

func main() {
	ensureIPLogFileExists()
	prevIp := getPreviousIp()
	currentIp := getCurrentIp()
	if prevIp == currentIp {
		return
	}
	updateCloudflareDNS(currentIp)
}

func ensureIPLogFileExists() {
	if _, err := os.Stat(IP_LOG_PATH); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(IP_LOG_PATH)
		defer file.Close()
		if err != nil {
			log.Fatalf("Failed to create %s: %s\n", IP_LOG_PATH, err.Error())
		}
	}
}

func getPreviousIp() string {
	lastIp, err := exec.Command("tail", "-n1", IP_LOG_PATH).Output()

	if err != nil {
		log.Fatalf("Failed to read %s: %s\n", IP_LOG_PATH, err.Error())
	}

	return string(lastIp)
}

func saveCurrentIp(ipAddr string) {
	file, err := os.OpenFile(IP_LOG_PATH, os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s\n", IP_LOG_PATH, err.Error())
	}
	defer file.Close()
	file.WriteString(ipAddr)
	file.Sync()
}

func getCurrentIp() string {
	res, err := http.Get("https://cloudflare.com/cdn-cgi/trace")
	if err != nil {
		log.Fatalf("Error occured sending request to cloudflare: %s\n", err.Error())
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error parsing response body: %s\n", err.Error())
	}
	rawIp := strings.Split(string(body), "\n")[2]
	return strings.Split(rawIp, "=")[1]
}

func updateCloudflareDNS(ipAddr string) {

}
