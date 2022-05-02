package region

import (
	"encoding/binary"
	"fmt"
	"os"
	"unicode"

	"github.com/hhhzzzsss/build-sync-generator/util"
)

const temp_path string = "TEMP_REGION_DUMP"
const output_path string = "../plugins/BuildSync/REGION_DUMP"
const dim int = 256

type Region struct {
	ids     [dim][dim][dim]int
	palette []string
}

func (r *Region) AddPaletteBlock(block string) {
	r.palette = append(r.palette, block)
}

func (r *Region) Set(x, y, z, id int) {
	r.ids[y][z][x] = id
}

func (r *Region) Get(x, y, z int) int {
	return r.ids[y][z][x]
}

func (r *Region) ForEach(idGenerator func(x, y, z int) int) {
	var bar util.ProgressBar
	bar.Initialize(256)
	for y := 0; y < dim; y++ {
		for z := 0; z < dim; z++ {
			for x := 0; x < dim; x++ {
				r.ids[y][z][x] = idGenerator(x, y, z)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
}

func (r *Region) ForEachNormalized(idGenerator func(x, y, z float64) int) {
	var bar util.ProgressBar
	bar.Initialize(256)
	for y := 0; y < dim; y++ {
		for z := 0; z < dim; z++ {
			for x := 0; x < dim; x++ {
				xNorm := 2.0*float64(x)/float64(dim) - 1.0
				yNorm := 2.0*float64(y)/float64(dim) - 1.0
				zNorm := 2.0*float64(z)/float64(dim) - 1.0
				r.ids[y][z][x] = idGenerator(xNorm, yNorm, zNorm)
			}
		}
		bar.Play(y + 1)
	}
	bar.Finish()
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
	for y := 0; y < dim; y++ {
		for z := 0; z < dim; z++ {
			for x := 0; x < dim; x++ {
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
	for y := 0; y < dim; y++ {
		for z := 0; z < dim; z++ {
			for x := 0; x < dim; x++ {
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
