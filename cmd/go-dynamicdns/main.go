package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/Maxxxxxx-x/go-dynamicdns/config"
	"github.com/cloudflare/cloudflare-go"
)

const IP_LOG_PATH = "./previous-ips.log"

func main() {
	config := config.LoadConfig()
	ensureIPLogFileExists()
	prevIp := getPreviousIp()
	currentIp := getCurrentIp()
	if prevIp == currentIp {
        log.Printf("Current IP is the same as ireviously stored IP. Aborting...")
		return
	}
	saveCurrentIp(currentIp)

	api, err := cloudflare.NewWithAPIToken(config.Auth.Token)
	if err != nil {
		log.Fatalf("Failed to connect to cloudflare: %s\n", err.Error())
	}
	updateCloudflareDNS(config, api, currentIp)
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
	file, err := os.OpenFile(IP_LOG_PATH, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    defer file.Close()
	if err != nil {
		log.Fatalf("Failed to open file %s: %s\n", IP_LOG_PATH, err.Error())
	}

    _, err = file.WriteString(ipAddr)
    if err != nil {
        log.Fatalf("Failed to write IP to %s: %s\n", IP_LOG_PATH, err.Error())
    }

    file.Sync()

    log.Printf("Wrote current IP %s to %s\n", ipAddr, IP_LOG_PATH)
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

func getZoneId(api *cloudflare.API, domain string) (*cloudflare.ResourceContainer, error) {
	splitted := strings.Split(domain, ".")
	if len(splitted) < 2 {
		return nil, fmt.Errorf("%s did not contian a valid TLD", domain)
	}
	zoneName := strings.Join(splitted[len(splitted)-2:], ".")
	zoneId, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ZoneID from %s\n", zoneName)
	}
	return cloudflare.ZoneIdentifier(zoneId), nil
}

func getAllDNSRecord(ctx context.Context, api *cloudflare.API, zoneIdent *cloudflare.ResourceContainer, domain string) ([]cloudflare.DNSRecord, error) {
	dnsRecords, _, err := api.ListDNSRecords(ctx, zoneIdent, cloudflare.ListDNSRecordsParams{
		Type: "A",
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot locate DNS record for %s in Zone %s\n", domain, zoneIdent.Identifier)
	}
	return dnsRecords, nil
}

func updateCloudflareDNS(config config.Config, api *cloudflare.API, ipAddr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	zoneId, err := getZoneId(api, config.Domains[0])
	if err != nil {
		log.Fatal(err)
	}

	dnsRecords, err := getAllDNSRecord(ctx, api, zoneId, config.Domains[0])
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range dnsRecords {
        if !slices.Contains(config.Domains, record.Name) {
            log.Printf("Skipping %s. [Not included in config]\n", record.Name)
            continue
        }
        if record.Content == ipAddr {
            log.Printf("Skipping %s. [IP unchanged]\n", record.Name)
            continue
        }
        log.Printf("Updating IP Address for %s [From %s to %s]\n", record.Name, record.Content, ipAddr)
		_, err := api.UpdateDNSRecord(ctx, zoneId, cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Name:    record.Name,
			Type:    "A",
			Content: ipAddr,
		})
		if err != nil {
			log.Fatalf("Erorr occured while updating DNS record for %s: %s\n", record.Name, err.Error())
		}
	}
}
