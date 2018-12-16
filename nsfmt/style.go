// This file is part of Netsak
// Copyright (C) 2018 matterbeam
// Use of this source is governed by a GPLv3

package nsfmt

import (
	"bytes"
	"io"
)

// ColorMode identifies different color output modes.
type ColorMode int

const (
	// ColorModeOff disables any text (default set will be used).
	ColorModeOff ColorMode = iota
	// ColorMode4Bit specifies that only the basic 16 ANSI colors will be used.
	ColorMode4Bit
	// ColorMode8Bit specifies to use 256-color mode.
	ColorMode8Bit
	// ColorMode24Bit enables 24bit "truecolor" mode.
	ColorMode24Bit
)

// Style represents a complete text style, including both foreground
// and background color.  It is decoded in a 64-bit integer.
// The coding is as follows (MSB first):
//    [  flags  ][ reserved ][  attributes  ][  bg color  ][  fg color  ]
//     < 2 bit >  <  6 bit >  <    8 bit   >  <  24 bit  >  <  24 bit  >
//
//
type Style uint64

var attrs = map[Style][]byte{
	AttrReset:   []byte{'0', ';'},
	AttrBold:    []byte{'1', ';'},
	AttrDim:     []byte{'2', ';'},
	AttrUnder:   []byte{'4', ';'},
	AttrBlink:   []byte{'5', ';'},
	AttrReverse: []byte{'7', ';'},
	AttrDefFg:   []byte{'3', '9', ';'},
	AttrDefBg:   []byte{'4', '9', ';'},
}

// ANSI returns the ANSI escape sequence equivalent to the Style.
func (s Style) ANSI(mode ColorMode) string {
	var b bytes.Buffer

	b.Write([]byte{'\x1b', '['})

	if s&AttrReset != 0 {
		b.Write([]byte{'0', 'm'})
		return b.String()
	}

	// Check the text attributes.
	for att, val := range attrs {
		if s&att != 0 {
			b.Write(val)
		}
	}

	// Check the colors.
	if mode != ColorModeOff {
		if s&FlagFg != 0 && s&AttrDefFg == 0 {
			doAnsiColor(mode, 38, &b, int(s&BitMaskFg))
		}
		if s&FlagBg != 0 && s&AttrDefBg == 0 {
			doAnsiColor(mode, 48, &b, int(s&BitMaskBg))
		}
	}

	// Terminate the sequence.
	b.Truncate(b.Len() - 1)
	b.WriteByte('m')

	return b.String()
}

// RgbTo4Bit maps an RGB color to a standard ANSI 4bit color.
func RgbTo4Bit(r, g, b int) int {
	// The 16 standard color palette can be ordered using 3 bits (b, g, r) where each bit
	// represents the presence of a parameter.
	//
	// To simplify things we are going to normalize the presence with the average of the
	// three parameters, this means that a parameter will only be considered as "present"
	// if its value is greater or equal than the average of the three parameters.
	//
	// Then, we will check the value of the average to differenciate between standard
	// and high-intensity colors.

	if r == g && r == b {
		// There is a special catch with the white varieties (black, silver, grey and white)
		// since silver (192,192,192) is standard and grey (128,128,128) is the "bright black".
		switch {
		case r < 64:
			return 0 // Black
		case r <= 128:
			return 8 // Grey
		case r <= 192:
			return 7 // Silver
		default:
			return 15 // White
		}
	}

	var n int
	avg := (r + g + b) / 3

	if r >= avg {
		n |= 1
	}
	if g >= avg {
		n |= 2
	}
	if b >= avg {
		n |= 4
	}

	// If the average is greater or equal than 128, consider it as a bright color.
	if avg >= 128 {
		n += 8
	}

	return n
}

// clamp returns y if x < y, z if x > z and x if y <= x <= z.
func clamp(x, y, z int) int {
	if x < y {
		return y
	}
	if x > z {
		return z
	}
	return x
}

// RgbTo8Bit maps an RGB color to an 8-bit color.
func RgbTo8Bit(r, g, b int) int {
	if r == g && r == b {
		// The same amount of red, green and blue produce a variety white. This is from
		// bright white (255;255;255) to black (0;0;0) - no white at all.
		// We can then map all the white combinations to the grayscale palette:
		//     256 white combinations / 24 grayscale colors ~= 10
		// Since all the combinations don't fit in 10 equal groups, higher values will
		// overflow the last group, so we will include them in the 10th group like this:
		return clamp(232+r/10, 0, 255)
	}

	// The 256-color palette contains the 16 ANSI colors (0-15), 24 grayscale colors (232-255)
	// and 216 different colors (16-231) organized in the following grid:
	//     - 6 squares of 6x6 cells. Each square to the right increases the green parameter.
	//     - In each square, each cell to the right increases the blue parameter.
	//     - In each square, each row (from top to bottom) increases the red parameter.
	// The value range is the same for each parameter {0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	//
	// Since we have already taken care of the grayscale colors, we will map the rest of the
	// colors to this grid. We have to find the closest value in the grid for the RGB color
	// we have received.
	//
	// We will do something similar to the grayscale, dividing each parameter in 6 possible
	// values and combining them to get the resultant cell in the grid.

	// 256 possible values for each parameter / 6 possible values in the grid ~= 42
	r = clamp(r/42, 0, 5)
	g = clamp(g/42, 0, 5)
	b = clamp(b/42, 0, 5)
	return 16 + (r * 36) + (g * 6) + b
}

// These bit masks can be used to retreive specific attributes from a style.
const (
	BitMaskFg   = 0x0000000000ffffff // Foreground color.
	BitMaskBg   = 0x0000ffffff000000 // Background color.
	BitMaskAttr = 0x00ff000000000000 // Text attributes.
	BitMaskFlag = 0x3000000000000000 // Style flags.
)

// Attributes are not colors, but affect the way the text is displayed.
const (
	// AttrReset - all attributes off.
	AttrReset Style = 1 << (iota + 48)
	// AttrBold - bold or increased intensity.
	AttrBold
	// AttrDim - decreased intensity.
	AttrDim
	// AttrUnder - underlined text.
	AttrUnder
	// AttrBlink - slow blink, less than 150 per minute.
	AttrBlink
	// AttrReverse - swap background and foreground colors.
	AttrReverse
	// AttrDefFg - use the default foreground color.
	AttrDefFg
	// AttrDefBg - use the default background color.
	AttrDefBg
)

// Flags specify parameters that affect the Style.
const (
	// FlagFg specifies that the foreground color is present.
	FlagFg Style = 1 << (63 - iota)
	// FlagBg specifies that the background color is present.
	FlagBg
)

// Foreground colors.
const (
	// ColorDefault - default foreground color
	ColorDefault Style = Style(0)
	// ColorBlack - black foreground color
	ColorBlack Style = (FlagFg | Style(0))
	// ColorRed - red foreground color
	ColorRed Style = (FlagFg | Style(0x800000))
	// ColorGreen - green foreground color
	ColorGreen Style = (FlagFg | Style(0x008000))
	// ColorYellow - yellow foreground color
	ColorYellow Style = (FlagFg | Style(0x808000))
	// ColorBlue - blue foreground color
	ColorBlue Style = (FlagFg | Style(0x000080))
	// ColorMagenta - magenta foreground color
	ColorMagenta Style = (FlagFg | Style(0x800080))
	// ColorCyan - cyan foreground color
	ColorCyan Style = (FlagFg | Style(0x008080))
	// ColorWhite - white foreground color
	ColorWhite Style = (FlagFg | Style(0xc0c0c0))
	// ColorGray - light black foreground color
	ColorGray Style = (FlagFg | Style(0x808080))
	// ColorLRed - light red foreground color
	ColorLRed Style = (FlagFg | Style(0xff0000))
	// ColorLGreen - light green foreground color
	ColorLGreen Style = (FlagFg | Style(0x00ff00))
	// ColorLYellow - light yellow foreground color
	ColorLYellow Style = (FlagFg | Style(0xffff00))
	// ColorLBlue - light blue foreground color
	ColorLBlue Style = (FlagFg | Style(0x0000ff))
	// ColorLMagenta - light magenta foreground color
	ColorLMagenta Style = (FlagFg | Style(0xff00ff))
	// ColorLCyan - light cyan foreground color
	ColorLCyan Style = (FlagFg | Style(0x00ffff))
	// ColorLWhite - light white foreground color
	ColorLWhite Style = (FlagFg | Style(0xffffff))
)

// Background colors.
const (
	// ColorBBlack - black foreground color
	ColorBBlack Style = (FlagBg | Style(0))
	// ColorBRed - red foreground color
	ColorBRed Style = (FlagBg | Style(0x800000000000))
	// ColorBGreen - green foreground color
	ColorBGreen Style = (FlagBg | Style(0x008000000000))
	// ColorBYellow - yellow foreground color
	ColorBYellow Style = (FlagBg | Style(0x808000000000))
	// ColorBBlue - blue foreground color
	ColorBBlue Style = (FlagBg | Style(0x000080000000))
	// ColorBMagenta - magenta foreground color
	ColorBMagenta Style = (FlagBg | Style(0x800080000000))
	// ColorBCyan - cyan foreground color
	ColorBCyan Style = (FlagBg | Style(0x008080000000))
	// ColorBWhite - white foreground color
	ColorBWhite Style = (FlagBg | Style(0xc0c0c0000000))
	// ColorBGray - light black foreground color
	ColorBGray Style = (FlagBg | Style(0x808080000000))
	// ColorBLRed - light red foreground color
	ColorBLRed Style = (FlagBg | Style(0xff0000000000))
	// ColorBLGreen - light green foreground color
	ColorBLGreen Style = (FlagBg | Style(0x00ff00000000))
	// ColorBLYellow - light yellow foreground color
	ColorBLYellow Style = (FlagBg | Style(0xffff00000000))
	// ColorBLBlue - light blue foreground color
	ColorBLBlue Style = (FlagBg | Style(0x0000ff000000))
	// ColorBLMagenta - light magenta foreground color
	ColorBLMagenta Style = (FlagBg | Style(0xff00ff000000))
	// ColorBLCyan - light cyan foreground color
	ColorBLCyan Style = (FlagBg | Style(0x00ffff000000))
	// ColorBLWhite - light white foreground color
	ColorBLWhite Style = (FlagBg | Style(0xffffff000000))
)

// doAnsiColor returns a byte slice containing the sequence to express the
// supplied rgb color in the proper format according to the specified color mode.
func doAnsiColor(mode ColorMode, fgbg int, w io.Writer, rgb int) {
	b := make([]byte, 15)[:0]

	switch mode {
	case ColorMode4Bit:
		n := RgbTo4Bit((rgb>>16)&0xff, (rgb>>8)&0xff, (rgb & 0xff))
		if n > 7 {
			n = n - 8 + 60
		}
		b = append(append(b, Itoa(fgbg-8+n, 0)...), ';')
	case ColorMode8Bit:
		n := RgbTo8Bit((rgb>>16)&0xff, (rgb>>8)&0xff, (rgb & 0xff))
		b = append(b, Itoa(fgbg, 0)...)
		b = append(b, ';', '5', ';')
		b = append(append(b, Itoa(n, 0)...), ';')
	case ColorMode24Bit:
		b = append(b, Itoa(fgbg, 0)...)
		b = append(b, ';', '2', ';')
		b = append(append(b, Itoa((rgb>>16)&0xff, 0)...), ';')
		b = append(append(b, Itoa((rgb>>8)&0xff, 0)...), ';')
		b = append(append(b, Itoa(rgb&0xff, 0)...), ';')
	}

	w.Write(b)
}
