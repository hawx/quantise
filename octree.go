package quantise

import (
	"image"
	"image/color"
)

// returns the dth bit of n
func bit(n, d uint8) uint8 {
	if n&(1<<(7-d)) == 0 {
		return 0
	}
	return 1
}

type Oct struct {
	isLeaf bool

	// A node in an octree simply has eight children
	Children [8]*Oct

	// A leaf has a color, and count
	Color *color.Color
	Count uint64
}

func NewOctree() *Oct {
	return &Oct{isLeaf: false}
}

func (tree *Oct) justInsert(c *color.Color, r, g, b, depth uint8) {
	tree.Count += 1

	if tree.isLeaf {
		return
	}

	index := bit(r, depth)<<2 | bit(g, depth)<<1 | bit(b, depth)

	if tree.Children[index] == nil {
		if depth == 5 {
			tree.Children[index] = &Oct{isLeaf: true, Color: c, Count: 1}
			return
		}

		tree.Children[index] = &Oct{isLeaf: false, Count: 1}
	}

	tree.Children[index].justInsert(c, r, g, b, depth+1)
}

func (tree *Oct) children() []*Oct {
	nodes := []*Oct{}

	for i := 0; i < 8; i++ {
		child := tree.Children[i]
		if child != nil && !child.isLeaf {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (tree *Oct) deepest() []*Oct {
	nodes := []*Oct{}
	last := []*Oct{tree}

	for {
		for i := 0; i < len(last); i++ {
			nodes = append(nodes, last[i].children()...)
		}

		if len(nodes) == 0 {
			return last
		}

		last = nodes
		nodes = []*Oct{}
	}

	return last
}

func (tree *Oct) Leaves() []*Oct {
	if tree.isLeaf {
		return []*Oct{tree}
	}

	leaves := []*Oct{}
	for i := 0; i < 8; i++ {
		child := tree.Children[i]
		if child != nil {
			leaves = append(leaves, child.Leaves()...)
		}
	}

	return leaves
}

func (tree *Oct) Insert(c color.Color) {
	if len(tree.Leaves()) <= 256 {
		r, g, b, _ := c.RGBA()
		tree.justInsert(&c, uint8(r), uint8(g), uint8(b), 0)

	} else {
		deepest := tree.deepest()
		least := deepest[0]

		for _, node := range deepest {
			if node.Count < least.Count {
				least = node
			}
		}

		least.Average()
		tree.Insert(c)
	}
}

func (tree *Oct) average() (color.Color, uint64) {
	if tree == nil {
		return nil, 0
	}

	if tree.isLeaf {
		if tree.Color == nil {
			return nil, 0
		}
		return *tree.Color, tree.Count
	}

	var rt, gt, bt, ct uint64

	for i := 0; i < 8; i++ {
		child := tree.Children[i]

		avg, c := child.average()
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

func (tree *Oct) Average() {
	if tree.isLeaf {
		return
	}

	avg, ct := tree.average()

	tree.Children = [8]*Oct{}
	tree.isLeaf = true
	tree.Color = &avg
	tree.Count = ct
}

func (tree *Oct) Palette() color.Palette {
	colors := []color.Color{}
	leaves := tree.Leaves()

	for i := 0; i < len(leaves); i++ {
		colors = append(colors, *leaves[i].Color)
	}

	return colors
}

func QuantizeOctree(img image.Image) image.Image {
	octree := NewOctree()
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			octree.Insert(img.At(x, y))
		}
	}

	out := image.NewPaletted(bounds, octree.Palette())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			out.Set(x, y, c)
		}
	}

	return out
}
