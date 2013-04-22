package main

import (
	"math/rand"
	"strings"
	"bytes"
	"fmt"
)

const (
	NucleobaseChoices = "CATG"
)

type SNP struct {
	GeneID string `cf:"snps" key:"GeneID" cols:"Description,Value,Alleles"`
	Description string
	Value string
	Alleles string
}

func (snp SNP) String() string {
	return fmt.Sprintf("%d: %s (%s) - %s", snp.GeneID, snp.Value, snp.Alleles, snp.Description)
}

func genRandomString(length int) string {
	var buffer bytes.Buffer
	for i := 0; i <= length; i++ {
		buffer.WriteByte(byte(rand.Intn(26) + 65))
	}
	return buffer.String()
}

func genSNPValue() string {
	idx := rand.Intn(len(NucleobaseChoices))
	buf := strings.Split(NucleobaseChoices, "")
	base := buf[idx]
	return base
}

func genAlleles(base string) string {
	var alleleBuffer bytes.Buffer
	for _, v := range NucleobaseChoices {
		if base == string(v) || rand.Intn(2) > 0 {
			alleleBuffer.WriteString(string(v))
		}
	}
	if len(alleleBuffer.String()) == 0 {
		alleleBuffer.WriteString(genSNPValue())
	}
	return alleleBuffer.String()
}

func genSNP(keyPrefix string) *SNP {
	id := keyPrefix //+  genRandomString(rand.Intn(64))
	description := genRandomString(rand.Intn(2048))
	value := genSNPValue()
	alleles := genAlleles(value)
	datum := &SNP{id, description, value, alleles}
	return datum
}


