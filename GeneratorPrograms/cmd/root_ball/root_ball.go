package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/treegen"
)

func main() {
	var r region.Region

	r.AddPaletteBlock("air")
	r.AddPaletteBlock("oak_wood")

	fDim := float64(region.Dim)
	halfDim := fDim / 2

	roots := []*treegen.SkeletonNode{treegen.NewSkeletonNode(halfDim, halfDim, halfDim)}

	fmt.Println("Loading attractor sets...")
	numAttractors := 3000
	attractors := make([]*treegen.Attractor, 0, numAttractors)
	for len(attractors) < numAttractors {
		x := 2*rand.Float64() - 1
		y := 2*rand.Float64() - 1
		z := 2*rand.Float64() - 1
		cdist := math.Sqrt(x*x + y*y + z*z)
		if cdist <= 1 && (cdist >= 0.5 || rand.Float64() < 1.5) {
			ax := halfDim + 0.95*halfDim*x
			ay := halfDim + 0.95*halfDim*y
			az := halfDim + 0.95*halfDim*z
			attractors = append(attractors, treegen.NewAttractor(ax, ay, az))
		}
	}

	fmt.Println("Generating skeleton...")
	settings := treegen.GetDefaultSettings()
	settings.StepSize = 1
	settings.KillDistance = 2
	settings.AttractionRadius = 18
	settings.BaseThickness = 20
	skeleton := treegen.GenerateSkeleton(roots, attractors, settings)

	fmt.Println("Fleshing out branches from skeleton...")
	skeleton.ForEachNode(func(node *treegen.SkeletonNode) {
		x := node.GetDim(0)
		y := node.GetDim(1)
		z := node.GetDim(2)
		bx := int(x)
		by := int(y)
		bz := int(z)
		r.Set(bx, by, bz, 1)
		r.ForEachInSphere(x, y, z, node.GetThickness(), func(sx, sy, sz int, rad2 float64) {
			r.Set(sx, sy, sz, 1)
		})
	})

	r.Hollow()

	r.CountBlocks()
	r.CreateDump()
}
