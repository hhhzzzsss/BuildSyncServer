package main

import (
	"math"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/util"
)

func main() {
	var region region.Region

	region.AddPaletteBlock("air")
	region.AddPaletteBlock("cyan_concrete")
	region.AddPaletteBlock("light_blue_concrete")
	region.AddPaletteBlock("blue_concrete")

	region.ForEachNormalized(func(x, y, z float64) int {
		Z := util.Triplex{}
		C := util.MakeTriplex(x, y, z).Multiply(1.2)
		minXDist := 10.0
		minYDist := 10.0
		minZDist := 10.0
		for i := 0; i < 5; i++ {
			Z = Z.Pow(8).Add(C)
			if math.Abs(Z.X) < minXDist {
				minXDist = math.Abs(Z.X)
			}
			if math.Abs(Z.Y) < minYDist {
				minYDist = math.Abs(Z.Y)
			}
			if math.Abs(Z.Z) < minZDist {
				minZDist = math.Abs(Z.Z)
			}
			if Z.LengthSquared() >= 4.0 {
				return 0
			}
		}
		color := 1
		minDist := minXDist
		if minYDist < minDist {
			minDist = minYDist
			color = 2
		}
		if minZDist < minDist {
			minDist = minZDist
			color = 3
		}
		return color
	})

	region.Hollow()

	region.CountBlocks()
	region.CreateDump()
}
