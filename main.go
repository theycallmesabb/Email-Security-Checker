package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Result slice to store all domain checks
var Result []DomainInfo

// Struct to hold details of a domain
type DomainInfo struct {
	Domain      string `json:"domain"`
	HasMX       bool   `json:"hasMX"`
	HasSPF      bool   `json:"hasSPF"`
	SPFRecord   string `json:"spfRecord"`
	HasDMARC    bool   `json:"hasDMARC"`
	DMARCRecord string `json:"dmarcRecord"`
}

func main() {
	// Scanner to take user input from terminal
	router := gin.Default()
	router.GET("/check/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		info := checkDomain(domain)
		c.JSON(200, info)
	})

	router.Run(":8080")

	scanner := bufio.NewScanner(os.Stdin)

	// Print header in CSV style
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	// Loop: for every domain entered, check its records
	for scanner.Scan() {
		domain := scanner.Text()
		info := checkDomain(domain)   // check DNS records
		Result = append(Result, info) // add to result list
		printdomain(info)             // print result
	}

	// If there's a scanning error, show it
	if err := scanner.Err(); err != nil {
		log.Fatal("Error:", err)
	}
}

func runScanner() {
	scanner := bufio.NewScanner(os.Stdin) // Scanner for reading user input from terminal

	// Print CSV header
	fmt.Printf("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord\n")

	// Loop through each domain entered by user
	for scanner.Scan() {
		domain := scanner.Text()
		info := checkDomain(domain)       // Get domain info
		Result = append(Result, info)     // Store in results slice
		fmt.Printf("%v,%v,%v,%v,%v,%v\n", // Print output in CSV format
			info.Domain, info.HasMX, info.HasSPF, info.SPFRecord, info.HasDMARC, info.DMARCRecord)
	}

}

// Domain ka MX, SPF, aur DMARC records check karne wali function
func checkDomain(domain string) DomainInfo {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Print(err) // Agar error aayi toh log me print karega
	}
	if len(mxRecords) > 0 {
		hasMX = true // Agar MX record mila toh hasMX ko true kar diya
	}

	// TXT records check kar raha hai (SPF aur DMARC dono isi me aate hain)
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Print(err)
	}

	// SPF record dhund raha hai
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") { // Agar record "v=spf1" se shuru ho raha hai toh SPF valid hai
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// DMARC record check kar raha hai
	dmarRecord, err := net.LookupTXT("dmarc." + domain) // **NOTE:** "dmarc" ke baad dot (.) lagana zaroori hai!
	if err != nil {
		log.Print(err)
	}

	// DMARC record dhund raha hai
	for _, record := range dmarRecord {
		if strings.HasPrefix(record, "v=DMARC1") { // Agar "v=DMARC1" se shuru ho raha hai toh DMARC valid hai
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	// Return all info in a struct
	return DomainInfo{
		Domain:      domain,
		HasMX:       hasMX,
		HasSPF:      hasSPF,
		SPFRecord:   spfRecord,
		HasDMARC:    hasDMARC,
		DMARCRecord: dmarcRecord,
	}

}
func printdomain(info DomainInfo) {
	fmt.Printf("%v,%v,%v,%v,%v,%v\n",
		info.Domain, info.HasMX, info.HasSPF,
		info.SPFRecord, info.HasDMARC, info.DMARCRecord)
}
