package quantise

import (
	"image"
	"image/color"
)

type Strategy int

const (
	// Merge the colours representing the fewest pixels
	LEAST Strategy = iota

	// Merge the colours representing the greatest pixels
	MOST
)

type OctreeQuantiser struct {
	Size     int
	Depth    uint8
	Strategy Strategy
}

func (q OctreeQuantiser) Quantise(in image.Image) color.Palette {
	octree := &oct{isLeaf: false}
	bounds := in.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			octree.insert(in.At(x, y), q.Size, q.Depth, q.Strategy)
		}
	}

	return octree.palette()
}

// returns the dth bit of n
func bit(n, d uint8) uint8 {
	if n&(1<<(7-d)) == 0 {
		return 0
	}
	return 1
}

// An oct represents a node or leaf in the oct-tree.
type oct struct {
	isLeaf bool

	// A node in an octree simply has eight children
	children [8]*oct

	// A leaf has a color, and count
	color *color.Color
	count uint64
}

func (tree *oct) justInsert(c *color.Color, r, g, b, depth, maxDepth uint8) {
	tree.count += 1

	if tree.isLeaf {
		return
	}

	index := bit(r, depth)<<2 | bit(g, depth)<<1 | bit(b, depth)

	if tree.children[index] == nil {
		if depth == maxDepth {
			tree.children[index] = &oct{isLeaf: true, color: c, count: 1}
			return
		}

		tree.children[index] = &oct{isLeaf: false, count: 1}
	}

	tree.children[index].justInsert(c, r, g, b, depth+1, maxDepth)
}

func (tree *oct) deepest() []*oct {
	nodes := []*oct{}
	last := []*oct{tree}

	for {
		for _, l := range last {
			for _, child := range l.children {
				if child != nil && !child.isLeaf {
					nodes = append(nodes, child)
				}
			}
		}

		if len(nodes) == 0 {
			return last
		}

		last = nodes
		nodes = []*oct{}
	}
}

func (tree *oct) leaves() []*oct {
	if tree.isLeaf {
		return []*oct{tree}
	}

	leaves := []*oct{}
	for _, child := range tree.children {
		if child != nil {
			leaves = append(leaves, child.leaves()...)
		}
	}

	return leaves
}

func (tree *oct) insert(c color.Color, size int, maxDepth uint8, strategy Strategy) {
	if len(tree.leaves()) <= size {
		r, g, b, _ := c.RGBA()
		tree.justInsert(&c, uint8(r), uint8(g), uint8(b), 0, maxDepth)

	} else {
		deepest := tree.deepest()
		toMerge := deepest[0]

		for _, node := range deepest {
			if strategy == LEAST {
				if node.count < toMerge.count {
					toMerge = node
				}
			} else {
				if node.count > toMerge.count {
					toMerge = node
				}
			}
		}

		toMerge.average()
		tree.insert(c, size, maxDepth, strategy)
	}
}

func (tree *oct) average2() (color.Color, uint64) {
	if tree == nil {
		return nil, 0
	}

	if tree.isLeaf {
		if tree.color == nil {
			return nil, 0
		}
		return *tree.color, tree.count
	}

	var rt, gt, bt, ct uint64

	for i := 0; i < 8; i++ {
		child := tree.children[i]

		avg, c := child.average2()
		if avg == nil {
			continue
		}

		r, g, b, _ := avg.RGBA()
		rt += uint64(r) * c / 255
		gt += uint64(g) * c / 255
		bt += uint64(b) * c / 255
		ct += c
	}

	if ct == 0 {
		return nil, 0
	}

	return color.RGBA{uint8(rt / ct), uint8(gt / ct), uint8(bt / ct), 255}, ct
}

func (tree *oct) average() {
	if tree.isLeaf {
		return
	}

	avg, ct := tree.average2()

	tree.children = [8]*oct{}
	tree.isLeaf = true
	tree.color = &avg
	tree.count = ct
}

func (tree *oct) palette() color.Palette {
	colors := []color.Color{}
	leaves := tree.leaves()

	for i := 0; i < len(leaves); i++ {
		colors = append(colors, *leaves[i].color)
	}

	return colors
}
