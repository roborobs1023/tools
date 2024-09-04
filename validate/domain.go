package validate

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func CheckDomain(domain string) (bool, error) {
	var hasMX, hasSPF, hasDMARC bool
	var sprRecord, dmarcRecord string

	exists := false
	mxRecord, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error:%v\nDomain: %s", err, domain)

		if strings.Contains(err.Error(), "no such host") {
			return exists, err
		}
	}
	if len(mxRecord) > 0 {
		hasMX = true
		exists = true
	}
	txt_records, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error:%v\n", err)
		return exists, err
	}

	for _, record := range txt_records {
		if strings.HasPrefix(record, "v=spfi") {
			hasSPF = true
			sprRecord = record
			exists = true
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
		return exists, err
	}

	for _, records := range dmarcRecords {
		if strings.HasPrefix(records, "V=DMARC") {
			hasDMARC = true
			dmarcRecord = records
			exists = true
			break
		}
	}

	fmt.Printf("%v,%v,%v,%v,%v,%v", domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord)
	return exists, nil
}
