#crcbench

See the blog post here:
[The Search For a Faster CRC32](https://blog.fastmail.com/2015/12/03/the-search-for-a-faster-crc32/)

crcbench is a CRC32 benchmarking program written in Go. There are currently 3 benchmarks.

1. `benchmark1` has the allocation of the CRC function using a sync.Pool outside loop.
2. `benchmark2` has the allocation of the CRC function using a sync.Pool inside loop.
3. `benchmark2` has the allocation of the CRC function using a sync.Pool inside loop and does two writes per iteration to simulate a key/value database calculating the checksum for a key and value. The size of the writes for the key and value is half the total length.

This code is based on `bench.c` from [crc32-bench](https://github.com/robn/crc32-bench) but was written from scratch in Go.

	Usage of crcbench:
	  -a	run all benchmarks.
	  -b1
	    	benchmark 1.
	  -b2
	    	benchmark 2.
	  -b3
	    	benchmark 3.
	  -crc string
	    	ieee, castagnoli, koopman. (default "ieee")
	  -l int
	    	buffer length. (default 16)
	  -n int
	    	number of rounds. (default 100000000)
	  -ss
	    	standard sequence of benchmarks.
	  -z	all data is zero.


###Sample output from a MacBookPro Intel Core i7 2860QM 4 core @ 2.5 GHz

####Go 1.5.1 using IEEE CRC from standard library (ieee)


	leb@hula:~/gotest/src/leb.io/crcbench % crcbench -a -ss -z            
	benchamrk1: sync.Pool outside loop (ieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=92.38 secs, avg=92.38 msecs, rate=726.46 MB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=91.00 secs, avg=909.96 nsecs, rate=281.33 MB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=45.93 secs, avg=459.25 nsecs, rate=278.71 MB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=23.50 secs, avg=235.03 nsecs, rate=272.31 MB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=12.16 secs, avg=121.62 nsecs, rate=263.12 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=9.34 secs, avg=93.43 nsecs, rate=256.88 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=6.43 secs, avg=64.26 nsecs, rate=248.97 MB/sec
	benchmark2: sync.Pool inside loop (ieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=91.96 secs, avg=91.96 msecs, rate=729.75 MB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=96.03 secs, avg=960.32 nsecs, rate=266.58 MB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=51.69 secs, avg=516.94 nsecs, rate=247.61 MB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=28.98 secs, avg=289.84 nsecs, rate=220.81 MB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=17.76 secs, avg=177.62 nsecs, rate=180.16 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=14.79 secs, avg=147.90 nsecs, rate=162.27 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=12.10 secs, avg=121.05 nsecs, rate=132.18 MB/sec
	benchmark3: sync.Pool inside loop with 2 writes (ieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=91.71 secs, avg=91.71 msecs, rate=731.79 MB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=96.87 secs, avg=968.68 nsecs, rate=264.28 MB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=52.31 secs, avg=523.14 nsecs, rate=244.68 MB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=29.91 secs, avg=299.07 nsecs, rate=214.00 MB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=18.66 secs, avg=186.63 nsecs, rate=171.46 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=15.81 secs, avg=158.08 nsecs, rate=151.82 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=13.02 secs, avg=130.23 nsecs, rate=122.86 MB/sec
	leb@hula:~/gotest/src/leb.io/crcbench % 

####Go 1.5.1 using IEEE CRC "github.com/klauspost/crc32" (kieee)

	leb@hula:~/gotest/src/leb.io/crcbench % crcbench -a -ss -z -crc="kieee"
	benchamrk1: sync.Pool outside loop (kieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=24.92 secs, avg=24.92 msecs, rate=2.69 GB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=12.80 secs, avg=128.04 nsecs, rate=2.00 GB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=8.11 secs, avg=81.06 nsecs, rate=1.58 GB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=5.79 secs, avg=57.88 nsecs, rate=1.11 GB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=8.97 secs, avg=89.73 nsecs, rate=356.64 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=7.87 secs, avg=78.75 nsecs, rate=304.77 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=6.81 secs, avg=68.12 nsecs, rate=234.87 MB/sec
	benchmark2: sync.Pool inside loop (kieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=24.37 secs, avg=24.37 msecs, rate=2.75 GB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=17.32 secs, avg=173.24 nsecs, rate=1.48 GB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=13.16 secs, avg=131.62 nsecs, rate=972.52 MB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=11.00 secs, avg=109.97 nsecs, rate=581.95 MB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=14.44 secs, avg=144.44 nsecs, rate=221.55 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=13.46 secs, avg=134.56 nsecs, rate=178.36 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=12.27 secs, avg=122.69 nsecs, rate=130.41 MB/sec
	benchmark3: sync.Pool inside loop with 2 writes (kieee)
		64 MiB, 1 krounds:	crc=0xb2eb30ed, runtime=24.27 secs, avg=24.27 msecs, rate=2.77 GB/sec
		256 B, 100 Mrounds:	crc=0x0d968558, runtime=20.37 secs, avg=203.67 nsecs, rate=1.26 GB/sec
		128 B, 100 Mrounds:	crc=0xc2a8fa9d, runtime=16.40 secs, avg=163.96 nsecs, rate=780.68 MB/sec
		64 B, 100 Mrounds:	crc=0x758d6336, runtime=23.06 secs, avg=230.62 nsecs, rate=277.52 MB/sec
		32 B, 100 Mrounds:	crc=0x190a55ad, runtime=18.57 secs, avg=185.69 nsecs, rate=172.33 MB/sec
		24 B, 100 Mrounds:	crc=0xa3c1ca20, runtime=16.44 secs, avg=164.43 nsecs, rate=145.96 MB/sec
		16 B, 100 Mrounds:	crc=0xecbb4b55, runtime=13.55 secs, avg=135.48 nsecs, rate=118.10 MB/sec
	leb@hula:~/gotest/src/leb.io/crcbench % 
