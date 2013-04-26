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
readStartDelay int, keyspace string, columnfamily string) (int64, int64) {
	importSnps := make([]*SNP, recordCount)
	for i := 0; i < recordCount; i++ {
		importSnps[i] = genSNP("")
	}
	startWriteTime := time.Now().UTC()

	mapping, err := gossie.NewMapping(&SNP{})
	if err != nil {
		exitMsg(fmt.Sprint("Creating mapping - ", err))
	}
	testSnp := genSNP("hecuna")
	var snpRow *gossie.Row
	snpRow, err = mapping.Map(testSnp)
	if err != nil {
		exitMsg(fmt.Sprint("Mapping SNP - ", err))
	}
	err = pool.Writer().Insert(columnfamily, snpRow).Run()
	if err != nil {
		exitMsg(fmt.Sprint("Write - ", err))
	}
	endWriteTime := time.Now().UTC()
	nsWrite := deltaNanoSeconds(endWriteTime, startWriteTime)
	startReadTime := time.Now().UTC()
	// Read
	endReadTime := time.Now().UTC()
	// Read
	nsRead := deltaNanoSeconds(endReadTime, startReadTime)
	return nsRead, nsWrite
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

	tRead, tWrite := benchmark(pool, *rowCount, 0, *keySpace,
*colFamily)
	fmt.Printf("Read time: %d\nWrite time: %d\n", tRead, tWrite)

}
