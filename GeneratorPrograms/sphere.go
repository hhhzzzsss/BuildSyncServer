package main

import (
	"github.com/hhhzzzsss/build-sync-generator/region"
)

func main() {
	var region region.Region

	region.AddPaletteBlock("air")
	region.AddPaletteBlock("stone")

	region.ForEachNormalized(func(x, y, z float64) int {
		if x*x+y*y+z*z < 1.0 {
			return 1
		} else {
			return 0
		}
	})

	region.CreateDump()
}
