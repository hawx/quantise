# quantise

Provides a function to reduce the number of colours in an image.

``` golang
func doQuantising(in image.Image) image.Image {
  return quantise.Quantise(in, quantise.OctreeQuantiser{
    Depth: 6,
    Size: 50,
    Strategy: quantise.LEAST,
  })
}
```

- Depth gives the maximum depth the tree can reach, the deeper the tree the
  greater the detail can be captured, kinda anyway.

- Size is the number of colours to use in the final image.

- Strategy is either LEAST in which case the colours representing the fewest
  pixels are merged, or MOST which merges the colours representing the most.

The results aren't very accurate.


## command line

There is a simple command line wrapper available:

``` bash
$ go get hawx.me/code/quantise/cmd/quantise
$ quantise --depth 6 --size 128 < input.png > output.png
```
