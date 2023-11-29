package main

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

type IpInfo struct {
	StartIP   string `csv:"start_ip"`
	EndIp     string `csv:"end_ip"`
	JoinKey   string `csv:"join_key"`
	IsHosting bool   `csv:"hosting"`
	IsProxy   bool   `csv:"proxy"`
	IsTor     bool   `csv:"tor"`
	IsVPN     bool   `csv:"vpn"`
	IsRelay   bool   `csv:"relay"`
	IsService bool   `csv:"service"`
}

type DB struct {
	// Blocks contains all the IPs and their associated info
	blocks map[string]IpInfo

	// Trie containing IPv4 addresses
	trieIPv4 ipaddr.AddressTrie

	// Trie containing IPv6 addresses
	trieIPv6 ipaddr.AddressTrie
}

func (db *DB) ParseCSVFileDefault(filePath string) error {
	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}

	r := csv.NewReader(reader)

	// Read the first line since it's the header and ignore it
	_, err = r.Read()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Record is an array of strings,
		//it does not differentiate between header and data
		recordInfo := IpInfo{
			StartIP:   record[0],
			EndIp:     record[1],
			JoinKey:   record[2],
			IsHosting: db.StringToBool(record[3]),
			IsProxy:   db.StringToBool(record[4]),
			IsTor:     db.StringToBool(record[5]),
			IsVPN:     db.StringToBool(record[6]),
			IsRelay:   db.StringToBool(record[7]),
			IsService: db.StringToBool(record[8]),
		}

		// This will return an array containing CIDR blocks for the star and end ip
		cidrBlocks := toPrefixBlock(recordInfo.StartIP, recordInfo.EndIp)

		// These each of these blocks is then added to a map
		// blocks[block.String()] = recordInfo
		// but we're not going to do that
		for _, block := range cidrBlocks {
			cidrBlockString := block.String()

			// Construct the IP address again from the CIDR block
			cidrAddress := ipaddr.NewIPAddressString(cidrBlockString)
			switch cidrAddress.GetIPVersion() {
			// Add it to the appropriate trie
			case ipaddr.IPv4:
				db.trieIPv4.Add(block.ToAddressBase())
			case ipaddr.IPv6:
				db.trieIPv6.Add(block.ToAddressBase())
			}
			db.blocks[cidrBlockString] = recordInfo
		}
	}
	return nil
}

func (db *DB) StringToBool(value string) bool {
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		return false
	}
}

// toPrefixBlock converts an IP range into a CIDR block. Adapted from
// https://github.com/seancfoley/ipaddress-go/wiki/Code-Examples-3:-Subnetting-and-Other-Subnet-Operations#from-start-and-end-address-get-a-minimal-list-of-cidr-blocks-spanning-the-range
func toPrefixBlock(start string, end string) []*ipaddr.IPAddress {
	startIP, endIP := ipaddr.NewIPAddressString(start), ipaddr.NewIPAddressString(end)
	startIPAddr, endIPAddr := startIP.GetAddress(), endIP.GetAddress()

	rng := startIPAddr.SpanWithRange(endIPAddr)

	blocks := rng.SpanWithPrefixBlocks()

	return blocks
}
