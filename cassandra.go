package main

import (
	"time"
	"fmt"
	"log"
	"github.com/carloscm/gossie/src/gossie"
)

var _ = log.Println

type CassandraEngine struct {
	pool gossie.ConnectionPool
	mapping gossie.Mapping
	keyspace string
	columnfamily string
}

func NewCassandraEngine(hosts []string, poolSize int, keyspace string, columnfamily string) (*CassandraEngine) {
	poolOptions := gossie.PoolOptions{Size: poolSize, Timeout: 3000}
	pool, err := gossie.NewConnectionPool(hosts, keyspace, poolOptions)
	if err != nil {
		ExitMsg(fmt.Sprint("Connecting: ", err))
	}
	engine := new(CassandraEngine)
	engine.pool = pool
	mapping, err := gossie.NewMapping(&SNP{})
	if err != nil {
		ExitMsg(fmt.Sprint("Creating mapping - ", err))
	}
	engine.mapping = mapping
	engine.keyspace = keyspace
	engine.columnfamily = columnfamily
	return engine
}

func (e *CassandraEngine) Benchmark(recordCount int) (BenchmarkData) {
	dataset := GenSNPDataset(recordCount)

	startWriteTime := time.Now().UTC()

	for _, snp := range dataset {
		snpRow, err := e.mapping.Map(snp)
		if err != nil {
			ExitMsg(fmt.Sprint("Mapping SNP - ", err))
		}
		mutation := e.pool.Writer().Insert(e.columnfamily, snpRow)
		err = mutation.Run()
		if err != nil {
			fmt.Println(fmt.Sprint("Write - ", err))
		}
	}
	endWriteTime := time.Now().UTC()
	startReadTime := time.Now().UTC()
	for _, snp := range dataset {
		query := e.pool.Query(e.mapping)
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
	summary := BenchmarkData{recordCount, startWriteTime, endWriteTime, startReadTime, endReadTime}
	return summary
}

