package main

/*

sortbyip

Sorts stdin by IP address.

Acceptable input lines:

    IP
    CIDR
    IP <whitespace> text
    CIDR <whitespace> text

  * IP can be an IPv4 or IPv6 address, or an IPv4 or IPv6 CIDR.
  * Whitespace is trimmed front and back.
  * The sort is stable.
  * Sorting two CIDRs with the same network portion but different mask lengths will be sorted sortest to longest network length.

*/

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strings"
)

type IpInfo struct {
	ip   net.IP
	line string
}

func readData(f io.Reader, ipList []IpInfo) ([]IpInfo, error) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse the line. It should start with an IP address.
		trim := strings.TrimSpace(line)
		i := strings.IndexAny(trim, "/ \t")
		if i == 0 { // Skip empty lines.
			continue
		}
		if i == -1 {
			i = len(trim)
		}
		if i < 4 { // Too short to be an IP? error
			return nil, fmt.Errorf("Line has no IP or CIDR: %v", line)
		}
		ip := net.ParseIP(trim[0:i])
		if ip == nil {
			return nil, fmt.Errorf("Line does not start with valid IP: %v", line)
		}

		ipList = append(ipList, IpInfo{ip, line})
	}

	return ipList, nil
}

func main() {
	flag.Parse()
	var ipList []IpInfo
	var err error

	// Read all input:
	if flag.NArg() == 0 {
		ipList, err = readData(os.Stdin, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, fname := range os.Args[1:] {
			fh, err := os.Open(fname)
			defer fh.Close()
			if err != nil {
				log.Fatal(err)
			}

			ipList, err = readData(fh, ipList)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Sort it.
	sort.SliceStable(ipList[:], func(i, j int) bool {
		return bytes.Compare(ipList[i].ip, ipList[j].ip) < 0
	})

	// Print it:
	for _, x := range ipList {
		fmt.Printf("%s\n", x.line)
	}

}
