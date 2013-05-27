package main

import (
	"math/rand"
	"fmt"
)

/*const (
	NucleobaseChoices = [...]byte {\C, \A, \T, \G}
)*/

var NucleobaseChoices = []byte("CATG")

type SNP struct {
	GeneID string `cf:"snps" key:"GeneID" cols:"Description"`
	Description string
	Value string
	Alleles string
}

func (snp SNP) String() string {
	return fmt.Sprintf("%d: %s (%s) - %s", snp.GeneID, snp.Value, snp.Alleles, snp.Description)
}

func genRandomString(length int) string {
	if length < 5 {
		length = 5
	}
	var buffer = make([]byte, length)
	for i := 0; i < length; i++ {
		asciiVal := rand.Intn(26) + 65
		if asciiVal < 65 || asciiVal > 90 {
			panic(fmt.Sprintf("%v", asciiVal))
		}
		buffer = append(buffer, byte(asciiVal))
	}
	str := string(buffer)
	return str
}

func genSNPValue() byte {
	idx := rand.Intn(len(NucleobaseChoices))
	base := NucleobaseChoices[idx]
	return base
}

func genAlleles(base byte) string {
	var alleleBuffer = make([]byte, 0)
	for _, v := range NucleobaseChoices {
		if base == v || rand.Intn(2) > 0 {
			alleleBuffer = append(alleleBuffer, v)
		}
	}
	if len(alleleBuffer) == 0 {
		alleleBuffer = append(alleleBuffer, genSNPValue())
	}
	return string(alleleBuffer)
}

func genSNP() *SNP {
	id := genRandomString(rand.Intn(64))
	description := genRandomString(rand.Intn(2048))
	value := genSNPValue()
	alleles := genAlleles(value)
	datum := &SNP{id, description, string(value), alleles}
	return datum
}


func GenSNPDataset(n int) ([]*SNP) {
	dataset := make([]*SNP, n)
        for i := 0; i < n; i++ {
                dataset[i] = genSNP()
	}
	return dataset
}
