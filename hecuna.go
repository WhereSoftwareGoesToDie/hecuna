/*
hecuna - Cassandra testing and benchmarking tool

Usage:

hecuna [options]

Options:

	keyspace (hecunatest)
	columnfamily (snps)
	rowcount (1000)
	hostlist (localhost:9160) (comma-separated host:port pairs)
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


func benchmarkCassandra(hosts []string, recordCount int,
poolSize int, keyspace string, columnfamily string) (int64, int64, int64, int64) {
	poolOptions := gossie.PoolOptions{Size: poolSize, Timeout: 3000}
	pool, err := gossie.NewConnectionPool(hosts, keyspace, poolOptions)
	if err != nil {
		exitMsg(fmt.Sprint("Connecting: ", err))
	}

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
	//	pool.Writer().Delete(columnfamily, []byte(snp.GeneID))
		mutation := pool.Writer().Insert(columnfamily, snpRow)
		err = mutation.Run()
		if err != nil {
			fmt.Println(fmt.Sprint("Write - ", err))
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
			fmt.Println(fmt.Sprint("Read - ", readErr))
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

type BackendFunc func([]string, int, int, string, string)(int64, int64, int64, int64)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	backends := map[string]BackendFunc{
		"cassandra": benchmarkCassandra,
	}

	hostList := flag.String("hosts", "localhost:9160","Comma-separated list of host:port pairs, e.g., localhost:9160,otherhost:9161")
	rowCount := flag.Int("rowcount", 1000, "Number of rows to write")
	keySpace := flag.String("keyspace", "hecunatest", "Name of keyspace to write to")
	colFamily := flag.String("colfamily", "snps", "Name of column family to write to")
	poolSize := flag.Int("poolsize", 50, "Number of connections in connection pool")

	flag.Usage = func() {
		backendList := ""
		for k := range backends {
			backendList += fmt.Sprintf("\t%s\n", k)
		}
		helpMessage := "hecuna: bringing clustered Java apps to their knees\n" +
			fmt.Sprintf("Usage: %s <backend> [options]\n\n", os.Args[0]) +
			"Backends:\n" +
			backendList +
			"\n" +
			"Options:\n\n"
		fmt.Fprintf(os.Stderr, helpMessage)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	backend := flag.Arg(0)

	hosts := strings.Split(*hostList, ",")

	tRead, tWrite, avgRead, avgWrite := backends[backend](hosts, *rowCount, *poolSize, *keySpace,
*colFamily)
	avgWriteSeconds := float64(avgWrite) / float64(1000000000)
	avgReadSeconds := float64(avgRead) / float64(1000000000)
	fmt.Printf("Read time: %v\n\tAverage (seconds): %v\nWrite time: %v\n\tAverage (seconds): %v\n", tRead, avgReadSeconds, tWrite, avgWriteSeconds)

}
