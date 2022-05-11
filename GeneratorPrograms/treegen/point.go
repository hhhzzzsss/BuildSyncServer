package treegen

import "github.com/hhhzzzsss/build-sync-generator/util"

type Point interface {
	GetDim(int) float64
}

func PointDistSq(a Point, b Point) float64 {
	dx := a.GetDim(0) - b.GetDim(0)
	dy := a.GetDim(1) - b.GetDim(1)
	dz := a.GetDim(2) - b.GetDim(2)
	return dx*dx + dy*dy + dz*dz
}

func PointToVec3d(p Point) util.Vec3d {
	return util.MakeVec3d(p.GetDim(0), p.GetDim(1), p.GetDim(2))
}
