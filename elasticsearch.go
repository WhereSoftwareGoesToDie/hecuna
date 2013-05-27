package main

import (
	"time"
	"log"
	"fmt"
	"github.com/mattbaird/elastigo/core"
)

var _ = log.Println

type ElasticsearchEngine struct {
	index string
	datatype string
}

func NewElasticsearchEngine(index string, datatype string) (*ElasticsearchEngine) {
	engine := new(ElasticsearchEngine)
	engine.index = index
	engine.datatype = datatype
	return engine
}

func (e *ElasticsearchEngine) Benchmark(recordCount int) (BenchmarkData) {
	dataset := GenSNPDataset(recordCount)

	startWriteTime := time.Now().UTC()

	for _, snp := range dataset {
		_, _ = core.Index(true, e.index, e.datatype, snp.GeneID, snp)
	}
	endWriteTime := time.Now().UTC()
	startReadTime := time.Now().UTC()
	for _, snp := range dataset {
		qry := fmt.Sprintf("%v:%v", e.datatype, snp.Value)
		_, _ = core.SearchUri(e.index, e.datatype, qry, "")
	}
	endReadTime := time.Now().UTC()
	summary := BenchmarkData{recordCount, startWriteTime, endWriteTime, startReadTime, endReadTime}
	return summary
}


