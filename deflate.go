package main

import (
	"deflate/compress"
	"flag"
	"fmt"
	"io"
	"os"
)

func checkErr(err error, exitCode int) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: ")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCode)
	}
}

func main() {
	// Parse arguments
	inPath := flag.String("in", "stdin", "specify input file")
	outPath := flag.String("out", "stdout", "specify output file")
	isCompression := flag.Bool("c", false, "compress file instead of decompression")
	blockSize := flag.Int("bs", 65536, "maximal block size in symbols")
	bufferInSize := flag.Int("insize", 258, "size of the input part of the sliding window (0-258)")
	bufferOutSize := flag.Int("outsize", 32768, "size of the output part of the sliding window (0-32768)")
	staticBlockThreshold := flag.Int("sthreshold", 256, "maximal block size in symbols that is encoded using static huffman trees")
	flag.Parse()

	// Open streams
	var (
		streamIn  io.Reader
		streamOut io.Writer

		err error
	)

	if *inPath == "stdin" {
		streamIn = os.Stdin
	} else {
		streamIn, err = os.Open(*inPath)
		checkErr(err, 1)
	}

	if *outPath == "stdout" {
		streamOut = os.Stdout
	} else {
		streamOut, err = os.Create(*outPath)
		checkErr(err, 2)
	}

	// Compress/decompress
	if *isCompression {
		checkErr(compress.Compress(
			streamIn,
			streamOut,
			*blockSize,
			*bufferInSize,
			*bufferOutSize,
			*staticBlockThreshold,
		), 3)
	} else {
		checkErr(compress.Decompress(streamIn, streamOut), 3)
	}
}
