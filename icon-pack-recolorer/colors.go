package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var uniqueColors = []string{}

var palette = []string{
	// "#f2d5cf",
	"#eebebe",
	"#f4b8e4",
	"#ca9ee6",
	"#e78284",
	"#ea999c",
	"#ef9f76",
	"#e5c890",
	"#a6d189",
	"#81c8be",
	"#99d1db",
	"#85c1dc",
	"#8caaee",
	"#babbf1",
	"#c6d0f5",
	"#b5bfe2",
	"#a5adce",
	"#949cbb",
	"#838ba7",
	"#737994",
	"#626880",
	"#51576d",
	"#414559",
	"#303446",
	"#292c3c",
	"#232634",
}

var explicitMatch = map[string]string{
	"#fff":    "#c6d0f5",
	"#ffffff": "#c6d0f5",
	"#fafafa": "#c6d0f5",
	"#f9f9f9": "#c6d0f5",
	"#dfdfdf": "#c6d0f5",
	"#dcdcdc": "#c6d0f5",
	"#b5b5b5": "#626880",
	"#333":    "#232634",
	"#ccc":    "#626880",
}

func hexToRGB(hex string) (float64, float64, float64, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex code: %s", hex)
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return float64(r), float64(g), float64(b), nil
}

func srgbToLinear(c float64) float64 {
	c = c / 255.0
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func fLab(t float64) float64 {
	delta := 6.0 / 29.0
	if t > math.Pow(delta, 3) {
		return math.Pow(t, 1.0/3.0)
	}
	return t/(3*delta*delta) + 4.0/29.0
}

func xyzToLab(x, y, z float64) (float64, float64, float64) {
	xn := 95.047
	yn := 100.000
	zn := 108.883

	fx := fLab(x / xn)
	fy := fLab(y / yn)
	fz := fLab(z / zn)

	L := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)

	return L, a, b
}

func rgbToXYZ(r, g, b float64) (float64, float64, float64) {
	rLin := srgbToLinear(r)
	gLin := srgbToLinear(g)
	bLin := srgbToLinear(b)

	x := rLin*0.4124 + gLin*0.3576 + bLin*0.1805
	y := rLin*0.2126 + gLin*0.7152 + bLin*0.0722
	z := rLin*0.0193 + gLin*0.1192 + bLin*0.9505

	return x * 100, y * 100, z * 100
}

func hexToLab(hex string) ([3]float64, error) {
	r, g, b, err := hexToRGB(hex)
	if err != nil {
		return [3]float64{}, err
	}

	x, y, z := rgbToXYZ(r, g, b)
	L, a, b := xyzToLab(x, y, z)
	return [3]float64{L, a, b}, nil
}

func deltaE(lab1, lab2 [3]float64) float64 {
	dL := lab1[0] - lab2[0]
	da := lab1[1] - lab2[1]
	db := lab1[2] - lab2[2]

	return math.Sqrt(dL*dL + da*da + db*db)
}

func hexSimilarity(color1 string, color2 string) float64 {
	lab1, err := hexToLab(color1)
	if err != nil {
		return math.MaxFloat64
	}
	lab2, err := hexToLab(color2)
	if err != nil {
		return math.MaxFloat64
	}

	return deltaE(lab1, lab2)
}

func ClosestMatchInPalette(hex string) string {
	explicitMatch, ok := explicitMatch[hex]
	if ok {
		return explicitMatch
	}

	contrasts := map[string]float64{}
	for _, color := range palette {
		contrasts[color] = hexSimilarity(hex, color)
	}

	var lowestKey string
	first := true
	for k := range contrasts {
		if first {
			lowestKey = k
			first = false
		} else if contrasts[k] < contrasts[lowestKey] {
			lowestKey = k
		}
	}
	return lowestKey
}
