package main

import (
	"testing"
	"fmt"
)

func TestSNPGeneration(t *testing.T) {
	for i := 0; i < 5; i++ {
		snp := genSNP("this-is-a-prefix")
		fmt.Println(snp, "\n")
	}
}

