package ipcache

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/danni-popova/ipinfo/ipdetails"
)

func Test_AddAndLookup(t *testing.T) {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test we can read from Redis
	_, err := rdb.Ping(ctx).Result()
	require.NoError(t, err)

	cache := New(rdb)
	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.0",
		EndIp:     "1.0.0.1",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
	})
	require.NoError(t, err)

	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.2",
		EndIp:     "1.0.0.2",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
		IsRelay:   true,
	})
	require.NoError(t, err)

	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.3",
		EndIp:     "1.0.0.45",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
	})
	require.NoError(t, err)

	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.46",
		EndIp:     "1.0.0.46",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
	})
	require.NoError(t, err)

	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.47",
		EndIp:     "1.0.0.49",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
	})
	require.NoError(t, err)

	err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
		StartIP:   "1.0.0.57",
		EndIp:     "1.0.0.67",
		JoinKey:   "1.0.0.0",
		IsHosting: true,
	})
	require.NoError(t, err)

	// Should find a matching record
	result, err := cache.LookupIP(ctx, "1.0.0.1")
	require.NoError(t, err)
	require.Equal(t, true, result.IsHosting)

	// Should find a matching record
	result, err = cache.LookupIP(ctx, "1.0.0.4")
	require.NoError(t, err)
	require.Equal(t, true, result.IsHosting)
	require.Equal(t, false, result.IsVPN)

	// Should not find a matching record
	result, err = cache.LookupIP(ctx, "192.0.0.4")
	require.NoError(t, err)
	require.Nil(t, result)

	result, err = cache.LookupIP(ctx, "1.0.0.68")
	require.NoError(t, err)
	require.Nil(t, result)

	result, err = cache.LookupIP(ctx, "1.0.0.56")
	require.NoError(t, err)
	require.Nil(t, result)
}

func Test_Lookup(t *testing.T) {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test we can read from Redis
	_, err := rdb.Ping(ctx).Result()
	require.NoError(t, err)

	cache := New(rdb)
	// Should find a matching record
	start := time.Now()
	_, err = cache.LookupIP(ctx, "24.117.240.164")
	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	require.NoError(t, err)

}
