package main

import (
	"flag"

	"hawx.me/code/img/utils"
	"hawx.me/code/quantise"
)

var (
	depth    = flag.Uint("depth", 8, "maximum depth of the tree to build")
	size     = flag.Int("size", 64, "number of colours to use")
	strategy = flag.String("strategy", "LEAST", "LEAST or MOST merge strategy")
)

func main() {
	flag.Parse()

	img, data := utils.ReadStdin()

	s := quantise.LEAST
	if *strategy == "MOST" {
		s = quantise.MOST
	}

	img = quantise.Quantise(img, quantise.OctreeQuantiser{
		Depth:    uint8(*depth),
		Size:     *size,
		Strategy: s,
	})

	utils.WriteStdout(img, data)
}
