package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

var mustMap = map[string]string{
	"#fff":    "#c6d0f5",
	"#ffffff": "#c6d0f5",
	"#fafafa": "#c6d0f5",
	"#f9f9f9": "#c6d0f5",
	"#dfdfdf": "#c6d0f5",
	"#dcdcdc": "#c6d0f5",
	// "#b5b5b5": "#626880",
	// "#333":    "#232634",
	// "#ccc":    "#626880",
}

type HaldCLUT struct {
	Size  int
	Level int
	Data  [][]color.NRGBA
	Cache map[string]string
}

func LoadCLUT(path string) (*HaldCLUT, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	if h <= 0 || w%h != 0 {
		return nil, fmt.Errorf("invalid LUT dimensions %dx%d", w, h)
	}

	clut := &HaldCLUT{
		Size:  h,
		Level: 8,
		Data:  make([][]color.NRGBA, h*h),
		Cache: make(map[string]string),
	}

	for y := range clut.Size * clut.Size {
		clut.Data[y] = make([]color.NRGBA, clut.Size)
	}

	for y := range h {
		for x := range w {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			clut.Data[y][x] = c
		}
	}

	return clut, nil
}

func MustLoadCLUT(path string) *HaldCLUT {
	clut, err := LoadCLUT(path)
	if err != nil {
		panic(err)
	}
	return clut
}

func (c *HaldCLUT) LookupPos(r, g, b uint8) (x, y int) {
	cubeSize := c.Level * c.Level

	ri := int(r) * (cubeSize - 1) / 255
	gi := int(g) * (cubeSize - 1) / 255
	bi := int(b) * (cubeSize - 1) / 255

	x = (ri % cubeSize) + (gi%c.Level)*cubeSize
	y = (bi * c.Level) + (gi / c.Level)

	return
}

func (c *HaldCLUT) Apply(in color.NRGBA) color.NRGBA {
	x, y := c.LookupPos(in.R, in.G, in.B)

	if y >= len(c.Data) || x >= len(c.Data[y]) {
		return color.NRGBA{255, 0, 255, 255}

	}

	return c.Data[y][x]
}

func (c *HaldCLUT) ClosestMatch(hex string) string {
	val, ok := mustMap[hex]
	if ok {
		return val
	}

	if cached, ok := c.Cache[hex]; ok {
		return cached
	}

	nrgba, err := HexToNRGBA(hex)
	if err != nil {
		fmt.Println(err)
		nrgba = color.NRGBA{255, 0, 255, 255}
	}
	outNRGBA := c.Apply(nrgba)
	out := NRGBAtoHexRGB(outNRGBA)

	c.Cache[hex] = out

	return out
}
