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
	r.AddPaletteBlock("oak_leaves")

	fDim := float64(region.Dim)
	halfDim := fDim / 2

	roots := []*treegen.SkeletonNode{treegen.NewSkeletonNode(halfDim, -5, halfDim)}
	rootRoots := []*treegen.SkeletonNode{treegen.NewSkeletonNode(halfDim, 0, halfDim)}

	fmt.Println("Loading attractor sets...")
	attractors := make([]*treegen.Attractor, 0, 10000)
	for len(attractors) < 10000 {
		x := 2*rand.Float64() - 1
		y := 2*rand.Float64() - 1
		z := 2*rand.Float64() - 1
		hdist := math.Sqrt(x*x + z*z)
		envelope := (2.5 - y) / 3 * math.Sqrt(1-y*y*y*y)
		excludedEnvelope := 0.
		radicand := 1 - 2*(y+1)*(y+1)
		if radicand >= 0 {
			excludedEnvelope = math.Sqrt(radicand) / 3
		}
		y01 := (y + 1) / 2
		if hdist <= envelope && hdist >= excludedEnvelope {
			ax := halfDim + 0.95*halfDim*x
			ay := 0.2*fDim + 0.75*fDim*y01
			az := halfDim + 0.95*halfDim*z
			attractors = append(attractors, treegen.NewAttractor(ax, ay, az))
		}
	}
	rootAttractors := make([]*treegen.Attractor, 0, 1000)
	for len(rootAttractors) < 1000 {
		x := 2*rand.Float64() - 1
		z := 2*rand.Float64() - 1
		if x*x+z*z <= 1 {
			ax := halfDim + 60*x
			ay := 0.
			az := halfDim + 60*z
			rootAttractors = append(rootAttractors, treegen.NewAttractor(ax, ay, az))
		}
	}

	fmt.Println("Generating main tree skeleton...")
	settings := treegen.GetDefaultSettings()
	settings.StepSize = 1
	settings.KillDistance = 2
	settings.AttractionRadius = 100
	settings.BaseThickness = 15
	settings.BaseBulge = 5
	settings.TaperThreshold = 4
	skeleton := treegen.GenerateSkeleton(roots, attractors, settings)

	fmt.Println("Generating root skeleton...")
	rootSettings := treegen.GetDefaultSettings()
	rootSettings.StepSize = 1
	rootSettings.KillDistance = 5
	rootSettings.AttractionRadius = 10
	rootSettings.BranchDecay = 0.02
	rootSettings.BaseThickness = 25
	rootSkeleton := treegen.GenerateSkeleton(rootRoots, rootAttractors, rootSettings)

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
	fmt.Println("Fleshing out roots from skeleton...")
	rootSkeleton.ForEachNode(func(node *treegen.SkeletonNode) {
		if node.IsRoot() {
			return
		}
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
	fmt.Println("Generating leaf cache...")
	var leafProbCache region.RegionCache[float64]
	skeleton.ForEachNode(func(node *treegen.SkeletonNode) {
		x := node.GetDim(0)
		y := node.GetDim(1)
		z := node.GetDim(2)
		bx := int(x)
		by := int(y)
		bz := int(z)
		if node.GetThickness() < 0.3 {
			r.Set(bx, by, bz, 2)
		}
		if node.GetThickness() < 1. {
			r.ForEachInSphere(x, y, z, node.GetThickness()+3, func(sx, sy, sz int, rad2 float64) {
				dist := math.Sqrt(rad2) - node.GetThickness()
				prob := (3 - dist) / 3
				leafProbCache.Set(sx, sy, sz, math.Max(leafProbCache.Get(sx, sy, sz), prob))
			})
		}
	})
	fmt.Println("Creating leaves...")
	r.ForEach(func(x, y, z int) int {
		if r.Get(x, y, z) != 0 {
			return r.Get(x, y, z)
		}
		prob := leafProbCache.Get(x, y, z)
		if prob > 0 && rand.Float64() < prob {
			return 2
		} else {
			return r.Get(x, y, z)
		}
	})

	r.SelectiveHollow(1)

	r.CountBlocks()
	r.CreateDump()
}
