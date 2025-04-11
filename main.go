package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)                             // User input lene ke liye scanner banaya
	fmt.Printf("domain,hasMX,hasSPF,spRecord,hasDMARC,dmarcRecord\n") // CSV format me output dene ke liye header

	// Har inputted domain ka check karega
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	// Agar scanner me koi error aata hai toh usko handle karega
	if err := scanner.Err(); err != nil {
		log.Fatal("Error", err)
	}
}

// Domain ka MX, SPF, aur DMARC records check karne wali function
func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spRecord, dmarcRecord string

	// MX records check kar raha hai
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
			spRecord = record
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

	// Final result print kar raha hai CSV format me
	fmt.Printf("%v,%v,%v,%v,%v,%v\n", domain, hasMX, hasSPF, spRecord, hasDMARC, dmarcRecord)
}
