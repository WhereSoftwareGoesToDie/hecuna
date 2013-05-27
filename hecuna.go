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
)

var _ = log.Println

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	engines := []string {
		"cassandra",
		"elasticsearch",
	}

	hostList := flag.String("hosts", "localhost:9160","Comma-separated list of host:port pairs, e.g., localhost:9160,otherhost:9161")
	rowCount := flag.Int("rowcount", 1000, "Number of rows to write")
	keySpace := flag.String("cassandra-keyspace", "hecunatest", "Name of keyspace to write to")
	colFamily := flag.String("cassandra-colfamily", "snps", "Name of column family to write to")
	poolSize := flag.Int("cassandra-poolsize", 50, "Number of connections in connection pool")
	esIndex := flag.String("elasticsearch-index", "snps", "Name of Elasticsearch index to write to")
	esDatatype := flag.String("elasticsearch-datatype", "snp", "Name of Elasticsearch datatype")

	flag.Usage = func() {
		backendList := ""
		for _, k := range engines {
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

	var engine StorageEngine
	
	switch backend {
		case "cassandra":
			engine = NewCassandraEngine(hosts, *poolSize, *keySpace, *colFamily)
		case "elasticsearch":
			engine = NewElasticsearchEngine(*esIndex, *esDatatype)
	}

	data := engine.Benchmark(*rowCount)
	fmt.Println(data)

}
