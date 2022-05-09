package main

import (
	"math"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/util"
)

func main() {
	var reg region.Region

	transparent := true
	if transparent {
		reg.AddPaletteBlock("air")
		reg.AddPaletteBlock("black_stained_glass")
		reg.AddPaletteBlock("white_stained_glass")
		reg.AddPaletteBlock("purple_stained_glass")
	} else {
		reg.AddPaletteBlock("air")
		reg.AddPaletteBlock("smooth_basalt")
		reg.AddPaletteBlock("calcite")
		reg.AddPaletteBlock("amethyst_block")
	}

	reg.ForEachNormalized(func(x, y, z float64) int {
		if y > 0 {
			return 0
		}
		scale := 1.1
		resolution := 2. * scale / float64(region.Dim)
		x = x * scale
		y = y * scale
		z = z*scale + 0.2
		Z := util.MakeQuaternion(x, y, z, 0)
		C := util.MakeQuaternion(-2, 6, 15, -6).Scale(1. / 22.)
		dZLen := 1.
		ZLen2 := Z.LengthSquared()
		escapeTime := 256
		for i := 0; i < 256; i++ {
			dZLen *= 3 * Z.LengthSquared()
			Z = Z.Cube().Add(C)
			ZLen2 = Z.LengthSquared()
			if ZLen2 > 256. || dZLen > 1e50 {
				escapeTime = i
				break
			}
		}
		dist := 0.5 * math.Log(ZLen2) * math.Sqrt(ZLen2) / dZLen
		if dist < resolution {
			if escapeTime < 30 {
				return 1
			} else if escapeTime < 50 {
				return 2
			} else if escapeTime < 70 {
				return 3
			} else {
				return 0
			}
		} else {
			return 0
		}
	})

	reg.CountBlocks()
	reg.CreateDump()
}
