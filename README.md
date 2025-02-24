# Deflate
Simple implementation of the Deflate compression/decompression in Go.

## Building
### Go
There are no dependencies in this project, so you can build it directly using go:
```bash
$ go build .
```
This will create a `deflate` binary file in the project folder.

### Docker (Linux only)
Also, there is a Dockerfile for building the project inside of the Docker container. To use it just execute the `build.sh` shell script:
```bash
$ ./build.sh
```
This will create a `deflate` binary in the `build` folder.

## Usage
Just execute the `deflate` binary with the according parameters.

The program expects input at `stdin`, writes output to `stdout`, and shows errors in `stderr`. This behavior can be overwritten by `in` and `out` flags.

By default program decompresses the input stream to the output stream, to compress specify `c` flag. The compression rate can be regulated by `bs`, `insize`, `outsize`, and `sthreshold` flags.

More information is available at `deflate -h`:

```bash
$ deflate -h
Usage of ./deflate:
  -bs int
        maximal block size in symbols (default 65536)
  -c    compress file instead of decompression
  -in string
        specify input file (default "stdin")
  -insize int
        size of the input part of the sliding window (0-258) (default 258)
  -out string
        specify output file (default "stdout")
  -outsize int
        size of the output part of the sliding window (0-32768) (default 32768)
  -sthreshold int
        maximal block size in symbols that is encoded using static huffman trees (default 256)
```

## Compression tools comparison

| Set | Zopfli | ZLIB | Deflate |
| - | - | - | - |
| README.md | 88.3% | 93.3% | 92.5% |
| alice29.txt | 33.8% | 35.7% | 36.4% |
| alphabet.txt | 0.3% | 0.3% | 0.5% |
| asyoulik.txt | 37.0% | 39.0% | 40.1% |
| cp.html | 31.3% | 32.4% | 33.2% |
| fields.c | 26.9% | 28.0% | 28.4% |
| grammar.lsp | 31.7% | 32.8% | 33.3% |
| helloworld.txt | 91.3% | 117.4% | 91.3% |
| random.txt | 75.2% | 77.1% | 77.3% |
| sum | 30.3% | 33.7% | 35.1% |
| xargs.1 | 39.9% | 41.1% | 41.6% |

We can see that the program works well compared to other tools, and in some examples even outperforms _ZLIB_ implementation. Yet it is not very close to _Zopfli_ which uses iterative improvement of the compression, and on average performs a little worse than _ZLIB_.

There were no time-measuring tests taken because this implementation is far (really far) behind those used widely (see next chapter).

## Place for improvement

#### TL;DR
Areas for improving:
- Hash chains instead of brute force in LZ77;
- Bit buffering;
- Package-merge algorithm optimization.

#### Details

This implementation is extremely slow due to the brute force (trying every possible (distance, length) pair) algorithm used in LZ77, while other implementations use hash chains. Calling functions for each bit operation without buffering (`BitStream`) builds up large function overhead, so 32-bit or 64-bit buffering would improve the situation. The Package-Merge algorithm is also not implemented with performance in mind.


Yet the main goal of the project was not to make an applicable compression tool but rather to understand better the compression format itself.
