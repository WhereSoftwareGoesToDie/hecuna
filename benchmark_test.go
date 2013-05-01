package main

import (
	"testing"
	"fmt"
)

func TestBenchmarkData(t *testing.T) {
	bd := &BenchmarkData{BenchmarkTime(4), BenchmarkTime(2)}
	fmt.Print(bd)
}

