package main

import (
	"github.com/hhhzzzsss/build-sync-generator/mmd_loader"
)

func main() {
	r := mmd_loader.LoadDumpAsRegion()
	r.CountBlocks()
	r.CreateDump()
}
