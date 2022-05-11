package region

import "github.com/hhhzzzsss/build-sync-generator/util"

type RegionCache[T any] struct {
	contents [Dim][Dim][Dim]T
}

func (r *RegionCache[T]) Set(x, y, z int, val T) {
	if !IsInRange(x, y, z) {
		return
	}
	r.contents[y][z][x] = val
}

func (r *RegionCache[T]) Get(x, y, z int) T {
	if !IsInRange(x, y, z) {
		return *new(T)
	}
	return r.contents[y][z][x]
}

func (r *RegionCache[T]) ForEach(generator func(x, y, z int) T) {
	var bar util.ProgressBar
	bar.Initialize(Dim)
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				r.contents[y][z][x] = generator(x, y, z)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
}

func (r *RegionCache[T]) ForEachNormalized(generator func(x, y, z float64) T) {
	var bar util.ProgressBar
	bar.Initialize(Dim)
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				xNorm := 2.0*float64(x)/float64(Dim) - 1.0
				yNorm := 2.0*float64(y)/float64(Dim) - 1.0
				zNorm := 2.0*float64(z)/float64(Dim) - 1.0
				r.contents[y][z][x] = generator(xNorm, yNorm, zNorm)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
}
