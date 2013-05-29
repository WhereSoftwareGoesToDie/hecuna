package main

import (
	"time"
	"fmt"
)

func deltaNanoSeconds(x, y time.Time) int64 {
	return x.UnixNano() - y.UnixNano()
}

type BenchmarkTime int64

type BenchmarkData struct {
	RecordCount int
	StartWrite time.Time
	StartRead time.Time
	EndWrite time.Time
	EndRead time.Time
}

// Time-per-record in nanoseconds
func (t BenchmarkTime) RecordAverage(count int) BenchmarkTime {
	return BenchmarkTime(int64(t) / int64(count))
}

func (t BenchmarkTime) Seconds() float64 {
	return float64(t) / float64(1000000000)
}

func (bd BenchmarkData) TimeWrite() BenchmarkTime {
	return BenchmarkTime(deltaNanoSeconds(bd.EndWrite, bd.StartWrite))
}

func (bd BenchmarkData) TimeRead() BenchmarkTime {
	return BenchmarkTime(deltaNanoSeconds(bd.EndRead, bd.StartRead))
}

func (bd BenchmarkData) AvgRead() float64 {
	return bd.TimeRead().RecordAverage(bd.RecordCount).Seconds()
}


func (bd BenchmarkData) AvgWrite() float64 {
	return bd.TimeWrite().RecordAverage(bd.RecordCount).Seconds()
}

func (bd BenchmarkData) String() string {
	summary := fmt.Sprint("Record count: ", bd.RecordCount)
	summary += fmt.Sprint("\nWrite time: ", bd.TimeWrite())
	summary += fmt.Sprint("\nRead time: ", bd.TimeRead())
	summary += fmt.Sprint("\nSeconds per write: ", bd.AvgWrite())
	summary += fmt.Sprint("\nSeconds per read: ", bd.AvgRead())
	return summary
}

type StorageEngine interface {
	Benchmark(recordCount int) BenchmarkData
}

