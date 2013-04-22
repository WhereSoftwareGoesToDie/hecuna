package main

import (
	"fmt"
	"testing"
	"github.com/carloscm/gossie/src/gossie"
)

func TestGossie(t *testing.T) {
	var keyspace = "hecunatest"
	var colspace = "snps"
	pool, err := gossie.NewConnectionPool([]string{"localhost:9160"}, keyspace, gossie.PoolOptions{Size: 50, Timeout: 3000})

	var mapping, mappingErr = gossie.NewMapping(&SNP{})
	if mappingErr != nil {
		exitMsg(fmt.Sprint("mapping: ", mappingErr))
	}

	var snpObj = &SNP{12, "does nothing", "C", "CAT"}
	var testSNP, mapErr = mapping.Map(snpObj)
	if mapErr != nil {
		exitMsg("Can't map value.")
	}
	err = pool.Writer().Insert(colspace, testSNP).Run()
	if err != nil {
		fmt.Println("Failed.")
	}
	fmt.Println("Connected.")
	var query = pool.Query(mapping)
	var ret, readErr = query.Get(12)
	if readErr != nil {
		exitMsg(fmt.Sprint("Couldn't read: ", readErr))
	}
	for {
		res_snp := &SNP{}
		err := ret.Next(res_snp)
		if err != nil {
			break
		}
		fmt.Println(res_snp.Alleles)
	}
	// Output: CAT
}

