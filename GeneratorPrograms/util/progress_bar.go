package util

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	cur   int
	total int
}

func (bar *ProgressBar) Initialize(total int) {
	bar.cur = 0
	bar.total = total
	bar.Play(0)
}

func (bar *ProgressBar) Play(cur int) {
	bar.cur = cur
	percent := int((float32(bar.cur) / float32(bar.total)) * 100)
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", strings.Repeat("=", percent/2), percent, bar.cur, bar.total)
}

func (bar *ProgressBar) Finish() {
	if bar.cur != bar.total {
		bar.Play(bar.total)
	}
	fmt.Println()
}
