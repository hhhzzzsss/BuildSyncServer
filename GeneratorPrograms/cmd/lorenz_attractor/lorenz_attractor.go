package main

import (
	"fmt"
	"math"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/util"
)

func main() {
	var r region.Region

	r.AddPaletteBlock("air")
	r.AddPaletteBlock("blue_concrete")
	r.AddPaletteBlock("cyan_concrete")
	r.AddPaletteBlock("light_blue_concrete")
	r.AddPaletteBlock("lime_concrete")
	r.AddPaletteBlock("yellow_concrete")
	r.AddPaletteBlock("orange_concrete")
	r.AddPaletteBlock("red_concrete")
	numColors := 7

	numIterations := 200000

	minX := 0.
	maxX := 0.
	minY := 0.
	maxY := 0.
	minZ := 0.
	maxZ := 0.

	rho := 28.
	sigma := 10.
	beta := 8. / 3.
	stepSize := 0.0005
	pos := util.MakeVec3d(1, 1, 1)
	fmt.Println("Doing dry iterations")
	for i := 0; i < numIterations; i++ {
		dx := sigma * (pos.Y - pos.X)
		dy := pos.X*(rho-pos.Z) - pos.Y
		dz := pos.X*pos.Y - beta*pos.Z
		dpos := util.MakeVec3d(dx, dy, dz)
		pos = pos.Add(dpos.Scale(stepSize))
	}
	for i := 0; i < numIterations; i++ {
		dx := sigma * (pos.Y - pos.X)
		dy := pos.X*(rho-pos.Z) - pos.Y
		dz := pos.X*pos.Y - beta*pos.Z
		dpos := util.MakeVec3d(dx, dy, dz)
		pos = pos.Add(dpos.Scale(stepSize))
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.Z < minZ {
			minZ = pos.Z
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
		if pos.Z > maxZ {
			maxZ = pos.Z
		}
	}

	maxDist := math.Max(maxX-minX, math.Max(maxY-minY, maxZ-minZ))
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	centerZ := (minZ + maxZ) / 2
	minX = centerX - 0.55*maxDist
	minY = centerY - 0.55*maxDist
	minZ = centerZ - 0.55*maxDist
	maxX = centerX + 0.55*maxDist
	maxY = centerY + 0.55*maxDist
	maxZ = centerZ + 0.55*maxDist

	fmt.Println("Doing block-placing iterations")
	for i := 0; i < numIterations; i++ {
		dx := sigma * (pos.Y - pos.X)
		dy := pos.X*(rho-pos.Z) - pos.Y
		dz := pos.X*pos.Y - beta*pos.Z
		dpos := util.MakeVec3d(dx, dy, dz)
		pos = pos.Add(dpos.Scale(stepSize))
		nx := (pos.X - minX) / (maxX - minX)
		ny := (pos.Y - minY) / (maxY - minY)
		nz := (pos.Z - minZ) / (maxZ - minZ)
		bx := int(math.Round(nx * float64(region.Dim)))
		by := int(math.Round(ny * float64(region.Dim)))
		bz := int(math.Round(nz * float64(region.Dim)))
		color := 1 + i/(numIterations/numColors)
		if color > numColors {
			color = numColors
		}
		r.Set(by, bx, bz, color)
	}

	r.CountBlocks()
	r.CreateDump()
}
