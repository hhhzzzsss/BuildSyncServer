package region

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"unicode"

	"github.com/hhhzzzsss/build-sync-generator/util"
)

const temp_path string = "TEMP_REGION_DUMP"
const output_path string = "../plugins/BuildSync/REGION_DUMP"
const Dim int = 256

type Region struct {
	ids     [Dim][Dim][Dim]int
	palette []string
}

func (r *Region) AddPaletteBlock(block string) {
	r.palette = append(r.palette, block)
}

func (r *Region) Set(x, y, z, id int) {
	if !IsInRange(x, y, z) {
		return
	}
	r.ids[y][z][x] = id
}

func (r *Region) Get(x, y, z int) int {
	if !IsInRange(x, y, z) {
		return 0
	}
	return r.ids[y][z][x]
}

func (r *Region) ForEach(idGenerator func(x, y, z int) int) {
	var bar util.ProgressBar
	bar.Initialize(Dim)
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				r.ids[y][z][x] = idGenerator(x, y, z)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
}

func (r *Region) ForEachNormalized(idGenerator func(x, y, z float64) int) {
	var bar util.ProgressBar
	bar.Initialize(Dim)
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				xNorm := 2.0*float64(x)/float64(Dim) - 1.0
				yNorm := 2.0*float64(y)/float64(Dim) - 1.0
				zNorm := 2.0*float64(z)/float64(Dim) - 1.0
				r.ids[y][z][x] = idGenerator(xNorm, yNorm, zNorm)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
}

func (r *Region) ForEachInSphere(cx, cy, cz, radius float64, f func(x, y, z int, rad2 float64)) {
	x1 := int(math.Ceil(cx - radius))
	y1 := int(math.Ceil(cy - radius))
	z1 := int(math.Ceil(cz - radius))
	x2 := int(math.Floor(cx + radius))
	y2 := int(math.Floor(cy + radius))
	z2 := int(math.Floor(cz + radius))
	for by := y1; by <= y2; by++ {
		for bz := z1; bz <= z2; bz++ {
			for bx := x1; bx <= x2; bx++ {
				dx := float64(bx) + 0.5 - cx
				dy := float64(by) + 0.5 - cy
				dz := float64(bz) + 0.5 - cz
				rad2 := dx*dx + dy*dy + dz*dz
				if rad2 <= radius*radius {
					f(bx, by, bz, rad2)
				}
			}
		}
	}
}

func (r *Region) Hollow() {
	fmt.Println("Hollowing...")
	var isSurface [Dim][Dim][Dim]bool
	for y := 1; y < Dim-1; y++ {
		for z := 1; z < Dim-1; z++ {
			for x := 1; x < Dim-1; x++ {
				if r.ids[y][z][x+1] == 0 ||
					r.ids[y][z][x-1] == 0 ||
					r.ids[y][z+1][x] == 0 ||
					r.ids[y][z-1][x] == 0 ||
					r.ids[y+1][z][x] == 0 ||
					r.ids[y-1][z][x] == 0 {
					isSurface[y][z][x] = true
				}
			}
		}
	}
	for y := 1; y < Dim-1; y++ {
		for z := 1; z < Dim-1; z++ {
			for x := 1; x < Dim-1; x++ {
				if !isSurface[y][z][x] {
					r.ids[y][z][x] = 0
				}
			}
		}
	}
}

func (r *Region) SelectiveHollow(id int) {
	fmt.Println("Hollowing...")
	var isSurface [Dim][Dim][Dim]bool
	for y := 1; y < Dim-1; y++ {
		for z := 1; z < Dim-1; z++ {
			for x := 1; x < Dim-1; x++ {
				if r.ids[y][z][x] != id {
					continue
				}
				if r.ids[y][z][x+1] != id ||
					r.ids[y][z][x-1] != id ||
					r.ids[y][z+1][x] != id ||
					r.ids[y][z-1][x] != id ||
					r.ids[y+1][z][x] != id ||
					r.ids[y-1][z][x] != id {
					isSurface[y][z][x] = true
				}
			}
		}
	}
	for y := 1; y < Dim-1; y++ {
		for z := 1; z < Dim-1; z++ {
			for x := 1; x < Dim-1; x++ {
				if r.ids[y][z][x] != id {
					continue
				}
				if !isSurface[y][z][x] {
					r.ids[y][z][x] = 0
				}
			}
		}
	}
}

func (r *Region) CountBlocks() {
	fmt.Print("Counting blocks...")
	numBlocks := 0
	for y := 1; y < Dim-1; y++ {
		for z := 1; z < Dim-1; z++ {
			for x := 1; x < Dim-1; x++ {
				if r.ids[y][z][x] != 0 {
					numBlocks++
				}
			}
		}
	}
	fmt.Printf("\rCounted a total of %d non-air blocks in the region\n", numBlocks)
	fmt.Printf("This would take roughly %.3f days (%.3f hours) to build\n", float64(numBlocks)/15/60/60/24, float64(numBlocks)/15/60/60)
}

func (r *Region) CreateDump() {
	r.Validate()

	fmt.Println("Writing region file...")
	f, err := os.Create(temp_path)
	if err != nil {
		panic(err)
	}

	binary.Write(f, binary.BigEndian, uint32(len(r.palette)))
	for _, paletteStr := range r.palette {
		for _, c := range paletteStr {
			if c > unicode.MaxASCII {
				panic("Palette entry was not ascii")
			}
		}
		binary.Write(f, binary.BigEndian, uint32(len(paletteStr)))
		binary.Write(f, binary.BigEndian, []byte(paletteStr))
	}
	var dataBuffer [256 * 256 * 256 * 4]byte
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				bufferIdx := y*256*256*4 + z*256*4 + x*4
				binary.BigEndian.PutUint32(dataBuffer[bufferIdx:], uint32(r.ids[y][z][x]))
			}
		}
	}
	binary.Write(f, binary.BigEndian, dataBuffer)

	f.Close()
	os.Rename(temp_path, output_path)
	fmt.Println("Finished creating region file")
}

// Panics if region has invalid state
func (r *Region) Validate() {
	fmt.Println("Validating region...")
	for _, paletteStr := range r.palette {
		for _, c := range paletteStr {
			if c > unicode.MaxASCII {
				panic("Palette entry was not ascii")
			}
		}
	}
	for y := 0; y < Dim; y++ {
		for z := 0; z < Dim; z++ {
			for x := 0; x < Dim; x++ {
				if r.ids[y][z][x] < 0 {
					errorMsg := fmt.Sprintf("Block id (%d) was less than zero", r.ids[y][z][x])
					panic(errorMsg)
				}
				if r.ids[y][z][x] >= len(r.palette) {
					errorMsg := fmt.Sprintf("Block id (%d) was greater than or equal to palette length (%d)", r.ids[y][z][x], len(r.palette))
					panic(errorMsg)
				}
			}
		}
	}
}

func IsInRange(x, y, z int) bool {
	return x >= 0 && x < Dim && y >= 0 && y < Dim && z >= 0 && z < Dim
}
