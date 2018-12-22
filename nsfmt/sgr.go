// This file is part of Netsak
// Copyright (C) 2018 matterbeam
// Use of this source is governed by a GPLv3

package nsfmt

import (
	"bytes"
	"strconv"
)

// These bit masks can be used to retreive specific attributes from an SGR.
const (
	BitMaskFg   = 0x0000000000ffffff // Foreground color.
	BitMaskBg   = 0x0000ffffff000000 // Background color.
	BitMaskAttr = 0x00ff000000000000 // Text attributes.
	BitMaskFlag = 0x3000000000000000 // SGR flags.
)

// ColorMode represents the way a color SGR sequence can be expanded.  There
// are three ways (or modes) that a color SGR can be expanded:
//     - 4-bit mode: allows 16 different colors (8 standard and 8 bright).
//     - 8-bit mode: allows 256 different colors.
//     - 24-bit mode ("truecolor"): allows for 16M different colors.
type ColorMode int

const (
	// ColorModeNone won't decode the colors in the SGR.
	ColorModeNone ColorMode = iota
	// ColorMode4Bit will decode the colors in the SGR using 4-bit mode.
	ColorMode4Bit
	// ColorMode8Bit will decode the colors in the SGR using 8-bit mode.
	ColorMode8Bit
	// ColorMode24Bit will decode the colors in the SGR using 24-bit mode.
	ColorMode24Bit
)

// SGR sets display attributes.  These attributes can change the color
// and font of the text displayed on terminals that support it.
//
// It is encoded in a 64-bit unsigned integer as follows (MSB first):
//    [  flags  ][ reserved ][  attributes  ][  bg color  ][  fg color  ]
//     < 2 bit >  <  6 bit >  <    8 bit   >  <  24 bit  >  <  24 bit  >
//
// This coding allows to combine different SGR using logial operations
// such as AND (&) and OR (|) to check or set parameters.
type SGR uint64

// DefaultSGR - Use default attributes.
const DefaultSGR = SGR(0)

// Attributes are not colors, but affect the way the text is displayed.
const (
	// AttrReset - all attributes off.
	AttrReset SGR = 1 << (iota + 48)
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

// Flags specify parameters that affect the SGR.
const (
	// FlagFg specifies that the foreground color is present.
	FlagFg SGR = 1 << (63 - iota)
	// FlagBg specifies that the background color is present.
	FlagBg
)

// Foreground colors.
const (
	// ColorBlack - black foreground color.
	ColorBlack SGR = (FlagFg | SGR(0))
	// ColorRed - red foreground color.
	ColorRed SGR = (FlagFg | SGR(0x800000))
	// ColorGreen - green foreground color.
	ColorGreen SGR = (FlagFg | SGR(0x008000))
	// ColorYellow - yellow foreground color.
	ColorYellow SGR = (FlagFg | SGR(0x808000))
	// ColorBlue - blue foreground color.
	ColorBlue SGR = (FlagFg | SGR(0x000080))
	// ColorMagenta - magenta foreground color.
	ColorMagenta SGR = (FlagFg | SGR(0x800080))
	// ColorCyan - cyan foreground color.
	ColorCyan SGR = (FlagFg | SGR(0x008080))
	// ColorWhite - white foreground color.
	ColorWhite SGR = (FlagFg | SGR(0xc0c0c0))
	// ColorGray - light black foreground color.
	ColorGray SGR = (FlagFg | SGR(0x808080))
	// ColorLRed - light red foreground color.
	ColorLRed SGR = (FlagFg | SGR(0xff0000))
	// ColorLGreen - light green foreground color.
	ColorLGreen SGR = (FlagFg | SGR(0x00ff00))
	// ColorLYellow - light yellow foreground color.
	ColorLYellow SGR = (FlagFg | SGR(0xffff00))
	// ColorLBlue - light blue foreground color.
	ColorLBlue SGR = (FlagFg | SGR(0x0000ff))
	// ColorLMagenta - light magenta foreground color.
	ColorLMagenta SGR = (FlagFg | SGR(0xff00ff))
	// ColorLCyan - light cyan foreground color.
	ColorLCyan SGR = (FlagFg | SGR(0x00ffff))
	// ColorLWhite - light white foreground color.
	ColorLWhite SGR = (FlagFg | SGR(0xffffff))
)

// Background colors.
const (
	// ColorBBlack - black foreground color.
	ColorBBlack SGR = (FlagBg | SGR(0))
	// ColorBRed - red foreground color.
	ColorBRed SGR = (FlagBg | SGR(0x800000000000))
	// ColorBGreen - green foreground color.
	ColorBGreen SGR = (FlagBg | SGR(0x008000000000))
	// ColorBYellow - yellow foreground color.
	ColorBYellow SGR = (FlagBg | SGR(0x808000000000))
	// ColorBBlue - blue foreground color.
	ColorBBlue SGR = (FlagBg | SGR(0x000080000000))
	// ColorBMagenta - magenta foreground color.
	ColorBMagenta SGR = (FlagBg | SGR(0x800080000000))
	// ColorBCyan - cyan foreground color.
	ColorBCyan SGR = (FlagBg | SGR(0x008080000000))
	// ColorBWhite - white foreground color.
	ColorBWhite SGR = (FlagBg | SGR(0xc0c0c0000000))
	// ColorBGray - light black foreground color.
	ColorBGray SGR = (FlagBg | SGR(0x808080000000))
	// ColorBLRed - light red foreground color.
	ColorBLRed SGR = (FlagBg | SGR(0xff0000000000))
	// ColorBLGreen - light green foreground color.
	ColorBLGreen SGR = (FlagBg | SGR(0x00ff00000000))
	// ColorBLYellow - light yellow foreground color.
	ColorBLYellow SGR = (FlagBg | SGR(0xffff00000000))
	// ColorBLBlue - light blue foreground color.
	ColorBLBlue SGR = (FlagBg | SGR(0x0000ff000000))
	// ColorBLMagenta - light magenta foreground color.
	ColorBLMagenta SGR = (FlagBg | SGR(0xff00ff000000))
	// ColorBLCyan - light cyan foreground color.
	ColorBLCyan SGR = (FlagBg | SGR(0x00ffff000000))
	// ColorBLWhite - light white foreground color.
	ColorBLWhite SGR = (FlagBg | SGR(0xffffff000000))
)

// attrs maps each attribute with its equivalent sequence.
var attrs = map[SGR]string{
	AttrReset:   "0;",
	AttrBold:    "1;",
	AttrDim:     "2;",
	AttrUnder:   "4;",
	AttrBlink:   "5;",
	AttrReverse: "7;",
	AttrDefFg:   "39;",
	AttrDefBg:   "49;",
}

// Decode decodes the SGR, returning a string containing the equivalent
// ANSI escape sequence.  Colors will be decoded using the color mode
// supplied.
func (s SGR) Decode(mode ColorMode) string {
	var b bytes.Buffer

	if s&AttrReset != 0 {
		// This attribute disables all text attributes, therefore
		// there is no need to continue decoding.
		return "\x1b[0m"
	}

	b.WriteString("\x1b[")

	// Check which attributes are set.
	for attr, val := range attrs {
		if s&attr != 0 {
			b.WriteString(val)
		}
	}

	// Decode the colors.
	if mode != ColorModeNone {
		if s&FlagFg != 0 && s&AttrDefFg == 0 { // FlagFg set and AttrDefFg not set.
			decodeColor(mode, 38, &b, int(s&BitMaskFg))
		}
		if s&FlagBg != 0 && s&AttrDefBg == 0 { // FlagBg set and AttrDefBg not set.
			decodeColor(mode, 48, &b, int(s&BitMaskBg))
		}
	}

	// Terminate the sequence.
	b.Truncate(b.Len() - 1)
	b.WriteByte('m')

	return b.String()
}

func (s SGR) String() string {
	return s.Decode(ColorMode8Bit)
}

// decodeColor decodes rgb using the supplied color mode.  The resulting
// sequence will be written into b.  fgbg must be adjusted to 38 or 48
// for a foreground or background color respectively.
func decodeColor(mode ColorMode, fgbg int, bf *bytes.Buffer, rgb int) {
	r, g, b := (rgb>>16)&0xff, (rgb>>8)&0xff, (rgb & 0xff)
	switch mode {
	case ColorMode4Bit:
		n := rgbTo4Bit(r, g, b)
		if n > 7 {
			n = n - 8 + 60
		}
		bf.WriteString(strconv.Itoa(fgbg-8+n) + ";")
	case ColorMode8Bit:
		n := rgbTo8Bit(r, g, b)
		bf.WriteString(strconv.Itoa(fgbg) + ";5;" + strconv.Itoa(n) + ";")
	case ColorMode24Bit:
		bf.WriteString(strconv.Itoa(fgbg) + ";2;" + strconv.Itoa(r) + ";" +
			strconv.Itoa(g) + ";" + strconv.Itoa(b) + ";")
	}
}

// rgbTo4Bit maps an RGB color to a standard ANSI 4bit color.
func rgbTo4Bit(r, g, b int) int {
	// The 16 standard color palette can be ordered using 3 bits (b, g, r)
	// where each bit represents the presence of a parameter.
	//
	// To simplify things we are going to normalize the presence with the
	// average of the three parameters, this means that a parameter will
	// only be considered as "present" if its value is greater or equal
	// than the average of the three parameters.
	//
	// Then, we will check the value of the average to differenciate between
	// standard and high-intensity colors.

	if r == g && r == b {
		// There is a special catch with the white varieties (black, silver,
		// grey and white) since silver (192,192,192) is standard and grey
		// (128,128,128) is the "bright black".
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

// rgbTo8Bit maps an RGB color to an 8-bit color.
func rgbTo8Bit(r, g, b int) int {
	if r == g && r == b {
		// The same amount of red, green and blue produce a variety white.  This
		// is from bright white (255;255;255) to black (0;0;0) - no white at all.
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
