/*
hecuna - Cassandra testing and benchmarking tool

Usage:

hecuna <keyspace>
*/

package main

import (
	"time"
	"math/rand"
)

func deltaNanoSeconds(x, y time.Time) int64 {
	return x.UnixNano() - y.UnixNano()
}

/*func benchmark(pool gossie.ConnectionPool, recordCount int, readStartDelay int) {
	// Generate SNPs
	startWriteTime := time.Now().UTC()
	// Write
	endWriteTime := time.Now().UTC()
	nsWrite := deltaNanoSeconds(endWrite, startWrite)
	startReadTime := time.Now().UTC()
	// Read
	nsRead := deltaNanoSeconds(endRead, startRead)
*/

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	testSNPGeneration()
}
