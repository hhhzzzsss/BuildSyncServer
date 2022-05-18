package mmd_loader

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hhhzzzsss/build-sync-generator/region"
	"github.com/hhhzzzsss/build-sync-generator/treegen"
)

const dump_path string = "MMD_DUMP"

func loadDump() *region.RegionCache[*Color] {
	f, err := os.Open(dump_path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bufr := bufio.NewReader(f)
	var cache region.RegionCache[*Color]
	cache.ForEach(func(x, y, z int) *Color {
		color := &Color{}
		binary.Read(bufr, binary.BigEndian, color)
		return color
	})
	return &cache
}

type BlockColorData []struct {
	Name  string    `json:"name"`
	Color []float32 `json:"color"`
}

func getBlockColorKDTree() *treegen.KDTree {
	bytes, _ := ioutil.ReadFile("resources/blockColors.json")
	var blockColorData BlockColorData
	json.Unmarshal(bytes, &blockColorData)
	kdtree := treegen.MakeKDTree()
	for _, blockData := range blockColorData {
		blockColor := MakeBlockColor(blockData.Color[0]/255, blockData.Color[1]/255, blockData.Color[2]/255, blockData.Color[3]/255, blockData.Name)
		kdtree.Add(&blockColor)
	}
	return &kdtree
}

func LoadDumpAsRegion() region.Region {
	cache := loadDump()
	kdtree := getBlockColorKDTree()
	var region region.Region
	region.AddPaletteBlock("air")
	indexMap := make(map[string]int)
	region.ForEach(func(x, y, z int) int {
		if cache.Get(x, y, z).A < 0.0 {
			return 0
		}
		blockColor := kdtree.NearestNeighbor(cache.Get(x, y, z)).(*BlockColor)
		blockIdx := region.PaletteSize()
		if idx, ok := indexMap[blockColor.Block]; ok {
			blockIdx = idx
		} else {
			indexMap[blockColor.Block] = blockIdx
			region.AddPaletteBlock(blockColor.Block)
		}
		return blockIdx
	})
	return region
}
