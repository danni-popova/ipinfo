package ipcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/ipinfo/ipdetails"
)

const (
	IPRedisKey = "iplookup" // TODO: This will be ipv4 or ipv6 depending on the version
)

type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) AddIPRange(ctx context.Context, info *ipdetails.IpRangeDetails) error {
	result := c.rdb.ZAdd(ctx, IPRedisKey, redis.Z{
		Score:  ipdetails.IpToFloat(info.EndIp),
		Member: info,
		// TODO: create a struct to hold all the remaining info in the member part
	})

	return result.Err()
}

func (c *Cache) LookupIP(ctx context.Context, ip string) (*ipdetails.IpRangeDetails, error) {
	ipFloat := ipdetails.IpToFloat(ip)

	result := c.rdb.ZRangeByScore(ctx, IPRedisKey, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", ipFloat),
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	})
	if result.Err() != nil {
		return nil, errors.New("couldn't do redis lookup")
	}

	if len(result.Val()) < 1 {
		fmt.Println("no ip ranges contain the given IP")
		return nil, nil
	}

	var resultDetails ipdetails.IpRangeDetails
	err := json.Unmarshal([]byte(result.Val()[0]), &resultDetails)
	if err != nil {
		return nil, err
	}

	if ipdetails.IpToFloat(resultDetails.StartIP) > ipFloat {
		fmt.Println("not in range")
		return nil, nil
	}

	return &resultDetails, nil
}
