/*
hecuna - Cassandra testing and benchmarking tool

Usage:

hecuna <keyspace>
*/

package main

import (
	"time"
	"math/rand"
	"log"
	"os"
	"fmt"
	"flag"
	"strings"
	"github.com/carloscm/gossie/src/gossie"
)

func exitMsg(msg string) {
	log.Println("Fatal:", msg)
	os.Exit(1)
}

func deltaNanoSeconds(x, y time.Time) int64 {
	return x.UnixNano() - y.UnixNano()
}

func benchmark(pool gossie.ConnectionPool, recordCount int,
readStartDelay int, keyspace string, columnfamily string) (int64, int64, int64, int64) {
	importSnps := make([]*SNP, recordCount)
	for i := 0; i < recordCount; i++ {
		importSnps[i] = genSNP()
	}
	startWriteTime := time.Now().UTC()

	mapping, err := gossie.NewMapping(&SNP{})
	if err != nil {
		exitMsg(fmt.Sprint("Creating mapping - ", err))
	}
	for _, snp := range importSnps {
		var snpRow *gossie.Row
		snpRow, err = mapping.Map(snp)
		if err != nil {
			exitMsg(fmt.Sprint("Mapping SNP - ", err))
		}
		mutation := pool.Writer().Insert(columnfamily, snpRow)
		err = mutation.Run()
		if err != nil {
			exitMsg(fmt.Sprint("Write - ", err))
		}
	}
	endWriteTime := time.Now().UTC()
	nsWrite := deltaNanoSeconds(endWriteTime, startWriteTime)
	avgWrite := nsWrite / int64(recordCount)
	startReadTime := time.Now().UTC()
	for _, snp := range importSnps {
		query := pool.Query(mapping)
		result, readErr := query.Get(snp.GeneID)
		if readErr != nil {
			exitMsg(fmt.Sprint("Read - ", readErr))
		}
		for {
			readSnp := &SNP{}
			err := result.Next(readSnp)
			if err != nil {
				break
			}
		}
	}
	endReadTime := time.Now().UTC()
	nsRead := deltaNanoSeconds(endReadTime, startReadTime)
	avgRead := nsRead / int64(recordCount)
	return nsRead, nsWrite, avgRead, avgWrite
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	hostList := flag.String("hosts", "localhost:9160","Comma-separated list of host:port pairs, e.g., localhost:9160,otherhost:9161")
	rowCount := flag.Int("rowcount", 1000, "Number of rows to write")
	keySpace := flag.String("keyspace", "hecunatest", "Name of keyspace to write to")
	colFamily := flag.String("colfamily", "snps", "Name of column family to write to")

	flag.Parse()

	hosts := strings.Split(*hostList, ",")
	poolOptions := gossie.PoolOptions{Size: 50, Timeout: 3000}

	pool, err := gossie.NewConnectionPool(hosts, *keySpace, poolOptions)
	if err != nil {
		exitMsg(fmt.Sprint("Connecting: ", err))
	}

	tRead, tWrite, avgRead, avgWrite := benchmark(pool, *rowCount, 0, *keySpace,
*colFamily)
	avgWriteSeconds := float64(avgWrite) / float64(1000000000)
	avgReadSeconds := float64(avgRead) / float64(1000000000)
	fmt.Printf("Read time: %v\n\tAverage (seconds): %v\nWrite time: %v\n\tAverage (seconds): %v\n", tRead, avgReadSeconds, tWrite, avgWriteSeconds)

}
