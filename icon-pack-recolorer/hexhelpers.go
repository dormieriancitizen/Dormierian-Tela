package main

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"
)

func NRGBAtoHexRGB(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func HexToNRGBA(s string) (color.NRGBA, error) {
	var out color.NRGBA
	if s == "" {
		return out, fmt.Errorf("empty string")
	}
	if s[0] == '#' {
		s = s[1:]
	}
	s = strings.TrimSpace(s)

	switch len(s) {
	case 3:
		expanded := make([]byte, 6)
		for i := range 3 {
			c := s[i]
			expanded[i*2] = c
			expanded[i*2+1] = c
		}
		s = string(expanded)
	case 4:
		expanded := make([]byte, 8)
		for i := range 4 {
			c := s[i]
			expanded[i*2] = c
			expanded[i*2+1] = c
		}
		s = string(expanded)
	case 6:
	case 8:
	default:
		return out, fmt.Errorf("invalid length: %d (expected 3,4,6 or 8 hex digits)", len(s))
	}

	switch len(s) {
	case 6:
		b, err := hex.DecodeString(s)
		if err != nil {
			return out, fmt.Errorf("invalid hex: %w", err)
		}
		out.R = b[0]
		out.G = b[1]
		out.B = b[2]
		out.A = 0xFF
		return out, nil
	case 8:
		b, err := hex.DecodeString(s)
		if err != nil {
			return out, fmt.Errorf("invalid hex: %w", err)
		}
		out.R = b[0]
		out.G = b[1]
		out.B = b[2]
		out.A = b[3]
		return out, nil
	default:

		return out, fmt.Errorf("unexpected parsing state")
	}
}
