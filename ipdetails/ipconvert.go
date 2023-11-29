package ipdetails

import (
	"net"
	"strconv"
	"strings"
)

func IpToFloat(ipString string) float64 {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return 0
	}
	// Convert each octet to a float64 and concatenate them
	parts := strings.Split(ip.String(), ".")
	var result float64
	for i := 0; i < len(parts); i++ {
		octet, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return 0
		}
		result = result*256 + octet
	}

	return result
}
