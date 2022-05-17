package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/treegen"
)

func main() {
	// Toggle to enable or disable leaf structures
	leaves := true

	var r region.Region

	r.AddPaletteBlock("air")
	r.AddPaletteBlock("blue_ice")
	r.AddPaletteBlock("magma_block")
	r.AddPaletteBlock("white_stained_glass")
	r.AddPaletteBlock("red_stained_glass")
	r.AddPaletteBlock("cobblestone")

	fDim := float64(region.Dim)
	halfDim := fDim / 2

	normCoord := 1 / math.Sqrt(3)
	iceRoot := treegen.NewSkeletonNode(halfDim+halfDim*normCoord, halfDim+halfDim*normCoord, halfDim+halfDim*normCoord)
	fireRoot := treegen.NewSkeletonNode(halfDim-halfDim*normCoord, halfDim-halfDim*normCoord, halfDim-halfDim*normCoord)
	roots := []*treegen.SkeletonNode{iceRoot, fireRoot}

	fmt.Println("Loading attractor sets...")
	numAttractors := 5000
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

	fmt.Println("Generating fire and ice tree skeleton...")
	settings := treegen.GetDefaultSettings()
	settings.StepSize = 1
	settings.KillDistance = 2
	settings.AttractionRadius = 100
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
		var id int
		if node.GetRoot() == iceRoot {
			id = 1
		} else if node.GetRoot() == fireRoot {
			id = 2
		}
		r.Set(bx, by, bz, id)
		r.ForEachInSphere(x, y, z, node.GetThickness(), func(sx, sy, sz int, rad2 float64) {
			r.Set(sx, sy, sz, id)
		})
	})

	if leaves {
		fmt.Println("Generating leaf cache...")
		var iceLeafCache region.RegionCache[float64]
		var fireLeafCache region.RegionCache[float64]
		skeleton.ForEachNode(func(node *treegen.SkeletonNode) {
			x := node.GetDim(0)
			y := node.GetDim(1)
			z := node.GetDim(2)
			if node.GetThickness() < 1. {
				r.ForEachInSphere(x, y, z, node.GetThickness()+3, func(sx, sy, sz int, rad2 float64) {
					dist := math.Sqrt(rad2) - node.GetThickness()
					prob := (2 - dist) / 2

					if node.GetRoot() == iceRoot {
						iceLeafCache.Set(sx, sy, sz, math.Max(iceLeafCache.Get(sx, sy, sz), prob))
					} else if node.GetRoot() == fireRoot {
						fireLeafCache.Set(sx, sy, sz, math.Max(fireLeafCache.Get(sx, sy, sz), prob))
					}
				})
			}
		})

		fmt.Println("Creating leaves...")
		r.ForEach(func(x, y, z int) int {
			if r.Get(x, y, z) != 0 {
				return r.Get(x, y, z)
			}
			iceFlag := false
			fireFlag := false
			probIce := iceLeafCache.Get(x, y, z)
			probFire := fireLeafCache.Get(x, y, z)
			if probIce > 0 && rand.Float64() < probIce {
				iceFlag = true
			}
			if probFire > 0 && rand.Float64() < probFire {
				fireFlag = true
			}

			if iceFlag && fireFlag {
				return 5
			} else if iceFlag {
				return 3
			} else if fireFlag {
				return 4
			} else {
				return r.Get(x, y, z)
			}
		})
	}

	r.SelectiveHollow(1)
	r.SelectiveHollow(2)

	r.CountBlocks()
	r.CreateDump()
}
