package main

import (
	"testing"
	"fmt"
)

func TestSNPGeneration(t *testing.T) {
	for i := 0; i < 100; i++ {
		snp := genSNP()
		fmt.Println(snp, "\n")
	}
}

