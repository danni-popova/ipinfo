package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/ipinfo/ipcache"
	"github.com/danni-popova/ipinfo/ipdetails"
)

const (
	FullIPInfoFile   = "ipinfo_privacy.csv"
	SampleIPInfoFile = "ipinfo_privacy_sample.csv"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test we can read from Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	cache := ipcache.New(rdb)

	reader, err := os.Open(FullIPInfoFile)
	if err != nil {
		fmt.Println(err)
		return
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
			fmt.Println(err)
			return
		}

		err = cache.AddIPRange(ctx, &ipdetails.IpRangeDetails{
			StartIP:   record[0],
			EndIp:     record[1],
			JoinKey:   record[2],
			IsHosting: StringToBool(record[3]),
			IsProxy:   StringToBool(record[4]),
			IsTor:     StringToBool(record[5]),
			IsVPN:     StringToBool(record[6]),
			IsRelay:   StringToBool(record[7]),
			IsService: StringToBool(record[8]),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
