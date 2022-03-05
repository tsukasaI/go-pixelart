package filter

import (
	"image/color"
	"log"
	"sort"
)

var (
	CGAPalettes map[int]color.Palette
	CGA64Table  []uint32
)

func init() {
	CGAPalettes = make(map[int]color.Palette)

	CGAPalettes[2] = color.Palette{color.Black, color.White}

	CGAPalettes[4] = generatePalette([]uint32{0x0, 0x55FFFF, 0xFF55FF, 0xFFFFFF})

	CGAPalettes[16] = generatePalette([]uint32{0x0, 0xAA, 0xAA00, 0xAAAA, 0xAA0000, 0xAA00AA, 0xAA5500, 0xAAAAAA,
		0x555555, 0x5555FF, 0x55FF55, 0x55FFFF, 0xFF5555, 0xFF55FF, 0xFFFF55, 0xFFFFFF})

	CGAPalettes[64] = generatePalette(initCGA64Table())
}

func generatePalette(t []uint32) color.Palette {
	c := make([]color.Color, len(t))
	sort.Slice(t, func(i, j int) bool { return sortAsc(t, i, j) })
	for i, e := range t {
		c[i] = CreateColor(e)
	}
	return c
}

func CGA(n int, c color.Color) color.Color {
	p, ok := CGAPalettes[n]
	if !ok {
		log.Fatalf("CGA%d not supported\n", n)
	}
	return p.Convert(c)
}

func initCGA64Table() []uint32 {
	s := make([]uint32, 64)
	for i := 0; i < 64; i++ {
		s[i] = convertBits(uint32(i), 3)
	}
	return s
}

func convertBits(x uint32, m int) uint32 {
	v := convertRightBits(x, m) + convertLeftBits(x, m)
	v &= 0xFFFFFF
	return v
}

func convertRightBits(x uint32, m int) uint32 {
	var v uint32 = 0
	for i := 0; i < m; i++ {
		if ((x >> i) & 0x1) == 0x1 {
			v |= (0xAA << (i * 8))
		}
	}
	return v
}

func convertLeftBits(x uint32, m int) uint32 {
	var v uint32 = 0
	x = x >> m
	for i := 0; i < m; i++ {
		if ((x >> i) & 0x1) == 0x1 {
			v |= (0x55 << (i * 8))
		}
	}
	return v
}

func sortAsc(s []uint32, i, j int) bool {
	return s[i] < s[j]
}

func FastCGA64(c color.Color, l bool) color.Color {
	r, g, b, _ := c.RGBA()
	r &= 0xFF
	g &= 0xFF
	b &= 0xFF

	var m uint32 = 0x4
	if !l {
		m = 0x3
	}

	r = (m * r) / 0x100
	r &= 0x3
	g = (m * g) / 0x100
	g &= 0x3
	b = (m * b) / 0x100
	b &= 0x3

	var v uint32 = 0
	v |= (r << 4)
	v |= (g << 2)
	v |= b
	v &= 0x3F

	return CGAPalettes[64][v]
}
