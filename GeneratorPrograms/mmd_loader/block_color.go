package mmd_loader

type Color struct {
	R, G, B, A float32
}

type BlockColor struct {
	R, G, B, A float32
	Block      string
}

func NewColor(r, g, b, a float32) *Color {
	return &Color{r, g, b, a}
}

func MakeBlockColor(r, g, b, a float32, block string) BlockColor {
	return BlockColor{r, g, b, a, block}
}

func (color *Color) GetDim(dim int) float64 {
	if dim == 0 {
		return float64(color.R)
	} else if dim == 1 {
		return float64(color.G)
	} else {
		return float64(color.B)
	}
}

func (blockColor *BlockColor) GetDim(dim int) float64 {
	if dim == 0 {
		return float64(blockColor.R)
	} else if dim == 1 {
		return float64(blockColor.G)
	} else {
		return float64(blockColor.B)
	}
}

// func RGB2LUV(r, g, b float64) (l, u, v float64) {
// 	if r > 0.04045 {
// 		r = math.Pow(((r + 0.055) / 1.055), 2.4)
// 	} else {
// 		r = r / 12.92
// 	}
// 	if g > 0.04045 {
// 		g = math.Pow(((g + 0.055) / 1.055), 2.4)
// 	} else {
// 		g = g / 12.92
// 	}
// 	if b > 0.04045 {
// 		b = math.Pow(((b + 0.055) / 1.055), 2.4)
// 	} else {
// 		b = b / 12.92
// 	}

// 	r *= 100
// 	g *= 100
// 	b *= 100

// 	X := r*0.4124 + g*0.3576 + b*0.1805
// 	Y := r*0.2126 + g*0.7152 + b*0.0722
// 	Z := r*0.0193 + g*0.1192 + b*0.9505

// 	u = (4 * X) / (X + (15 * Y) + (3 * Z))
// 	v = (9 * Y) / (X + (15 * Y) + (3 * Z))

// 	Y = Y / 100
// 	if Y > 0.008856 {
// 		Y = math.Pow(Y, 1./3.)
// 	} else {
// 		Y = (7.787 * Y) + (16. / 116.)
// 	}

// 	ref_X, ref_Y, ref_Z := 95.047, 100.000, 108.883
// 	ref_u := (4 * ref_X) / (ref_X + (15 * ref_Y) + (3 * ref_Z))
// 	ref_v := (9 * ref_Y) / (ref_X + (15 * ref_Y) + (3 * ref_Z))

// 	CIE_L := (116 * Y) - 16.
// 	CIE_U := 13 * CIE_L * (u - ref_u)
// 	CIE_V := 13 * CIE_L * (v - ref_v)

// 	return CIE_L, CIE_U, CIE_V
// }
