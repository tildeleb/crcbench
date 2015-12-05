// © Copyright 2015 Lawrence E. Bakst All Rights Reserved
// Based on the code found here: https://github.com/robn/crc32-bench

// This program tests CRC implementations with various size buffers and number of passes.
// IEEE, Castagnoli, and Koopman polnomials are supported.
// Data can be all zeros or random. A standard sequence of benchmarks values are included
// or you can specify the length and/or passes via command line flags.
package main

import (
	"flag"
	"fmt"
	kcrc32 "github.com/klauspost/crc32"
	"hash"
	"hash/crc32"
	"math/rand"
	"sync"
	"time"

	"leb.io/hrff"
)

// flags
var crc = flag.String("crc", "ieee", "ieee, kieee, castagnoli, koopman.")
var length = flag.Int("l", 16, "buffer length.")
var n = flag.Int("n", 100*1000*1000, "number of rounds.")
var af = flag.Bool("a", false, "run all benchmarks.")
var ss = flag.Bool("ss", false, "standard sequence of benchmarks.")
var zf = flag.Bool("z", false, "all data is zero.")

func init() {
	flag.BoolVar(&Benchmarks[0].Run, "b1", false, "benchmark 1.")
	flag.BoolVar(&Benchmarks[1].Run, "b2", false, "benchmark 2.")
	flag.BoolVar(&Benchmarks[2].Run, "b3", false, "benchmark 3.")
}

// BenchmarkConfig specifies a length and number of passes "N" for a benchmark run.
type BenchmarkConfig struct {
	Length int // size of the crc buffer
	N      int // number of passes
}

// A Benchmark is a function that runs the benchmark, a header to print, and a bool if we should run it.
type Benchmark struct {
	Benchmark func(l, n int) (sum uint32, datalen, run, avg float64) // method value
	H1Header  string                                                 // print header before benchmark
	Run       bool                                                   // true if running this benchmark
}

// The State for a benchmark is a buffer and a sync.Pool for the crc.
type State struct {
	Buf       []byte
	CRC32Pool sync.Pool
}

var state *State = new(State)

// Benchmarks holds the various benchmarks that can be run.
var Benchmarks = []Benchmark{
	{state.Benchmark1, "benchamrk1: sync.Pool outside loop", false},
	{state.Benchmark2, "benchmark2: sync.Pool inside loop", false},
	{state.Benchmark3, "benchmark3: sync.Pool inside loop with 2 writes", false},
}

const M100 = 100 * 1000 * 1000
const Mi64 = 64 * 1024 * 1024

// StandardSequence contains the length and passes run for each benchmark.
var StandardSequence = []BenchmarkConfig{
	BenchmarkConfig{Length: Mi64, N: 1000},
	BenchmarkConfig{Length: 256, N: M100},
	BenchmarkConfig{Length: 128, N: M100},
	BenchmarkConfig{Length: 64, N: M100},
	BenchmarkConfig{Length: 32, N: M100},
	BenchmarkConfig{Length: 24, N: M100},
	BenchmarkConfig{Length: 16, N: M100},
}

var src = rand.NewSource(time.Now().UTC().UnixNano())
var r = rand.New(src)

// rbetween returns a random int [a, b]
func rbetween(a int, b int) int {
	if *zf {
		return 0
	}
	return r.Intn(b-a+1) + a
}

// GetCRC returns a sync.Pool for the specified crc.
func GetCRC(crc string) sync.Pool {
	switch crc {
	case "kieee":
		return sync.Pool{New: func() interface{} { return kcrc32.NewIEEE() }}
	case "ieee":
		return sync.Pool{New: func() interface{} { return crc32.NewIEEE() }}
	case "castagnoli":
		return sync.Pool{New: func() interface{} { return crc32.New(crc32.MakeTable(crc32.Castagnoli)) }}
	case "koopman":
		return sync.Pool{New: func() interface{} { return crc32.New(crc32.MakeTable(crc32.Koopman)) }}
	}
	panic("unknown crc: " + crc)
}

// H2Header prints the results for each benchmark and returns the datalen
func H2Header(length, n int) float64 {
	fmt.Printf("\t%H, %h:\t", hrff.Int{length, "B"}, hrff.Int{n, "rounds"})
	return float64(length * n)
}

func (s *State) FillBuffer(length int) {
	// initialize buffer with random data or zeros
	s.Buf = make([]byte, length, length)
	for k := range s.Buf {
		s.Buf[k] = byte(rbetween(0, 255))
	}
}

func (s *State) Benchmark1(length, n int) (sum uint32, datalen, run, avg float64) {
	datalen = H2Header(length, n)
	s.FillBuffer(length)
	crc := s.CRC32Pool.Get().(hash.Hash32)
	beg := time.Now()
	for i := 0; i < n; i++ {
		if _, err := crc.Write(s.Buf); err != nil {
			panic(err)
		}
		sum = crc.Sum32()
		crc.Reset()
	}
	end := time.Now()
	s.CRC32Pool.Put(crc)
	run = float64(end.Sub(beg).Seconds())
	avg = run / float64(n)
	return
}

func (s *State) Benchmark2(length, n int) (sum uint32, datalen, run, avg float64) {
	datalen = H2Header(length, n)
	s.FillBuffer(length)
	beg := time.Now()
	for i := 0; i < n; i++ {
		crc := s.CRC32Pool.Get().(hash.Hash32)
		if _, err := crc.Write(s.Buf); err != nil {
			panic(err)
		}
		sum = crc.Sum32()
		crc.Reset()
		s.CRC32Pool.Put(crc)
	}
	end := time.Now()
	run = float64(end.Sub(beg).Seconds())
	avg = run / float64(n)
	return
}

func (s *State) Benchmark3(length, n int) (sum uint32, datalen, run, avg float64) {
	datalen = H2Header(length, n)
	s.FillBuffer(length)
	key := s.Buf[0 : length/2]
	val := s.Buf[length/2 : length]

	beg := time.Now()
	for i := 0; i < n; i++ {
		crc := s.CRC32Pool.Get().(hash.Hash32)
		if _, err := crc.Write(key); err != nil {
			panic(err)
		}
		if _, err := crc.Write(val); err != nil {
			panic(err)
		}
		sum = crc.Sum32()
		crc.Reset()
		s.CRC32Pool.Put(crc)
	}
	end := time.Now()
	run = float64(end.Sub(beg).Seconds())
	avg = run / float64(n)
	return
}

func main() {
	var print = func(sum uint32, datalen, run, avg float64) {
		fmt.Printf("crc=%#08x, runtime=%0.2f secs, avg=%0.2h, rate=%0.2h\n",
			sum, run, hrff.Float64{avg, "secs"}, hrff.Float64{datalen / run, "B/sec"})
	}
	flag.Parse()
	if *af {
		for k := range Benchmarks {
			Benchmarks[k].Run = true
		}
	}

	for _, v := range Benchmarks {
		if !v.Run {
			continue
		}
		fmt.Printf("%s (%s)\n", v.H1Header, *crc)
		state.CRC32Pool = GetCRC(*crc)
		benchmark := v.Benchmark
		if *ss {
			for _, k := range StandardSequence {
				print(benchmark(k.Length, k.N))
			}
		} else {
			print(benchmark(*length, *n))
		}
	}
}
