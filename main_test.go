package main

import (
	"testing"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/stretchr/testify/require"
)

func TestDB_ParseCSVFileDefault(t *testing.T) {
	ipInfoDB := &DB{
		blocks:   map[string]IpInfo{},
		trieIPv4: ipaddr.AddressTrie{},
		trieIPv6: ipaddr.AddressTrie{},
	}

	err := ipInfoDB.ParseCSVFileDefault(FullIPInfoFile)
	require.NoError(t, err)
}
