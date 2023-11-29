package ipcache

import (
	"context"
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
	startIp := ipdetails.IpToFloat(info.StartIP)

	result := c.rdb.ZAdd(ctx, IPRedisKey, redis.Z{
		Score:  ipdetails.IpToFloat(info.EndIp),
		Member: fmt.Sprintf("%f", startIp),
		// TODO: create a struct to hold all the remaining info in the member part
	})

	return result.Err()
}

func (c *Cache) LookupIP(ctx context.Context, ip string) error {
	ipFloat := ipdetails.IpToFloat(ip)

	result := c.rdb.ZRangeByScore(ctx, IPRedisKey, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", ipFloat),
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	})

	if len(result.Val()) < 1 {
		fmt.Println("no ip ranges contain the given IP")
		return nil
	}

	details := result.Val()[0]
	fmt.Println(details)

	return result.Err()
}
