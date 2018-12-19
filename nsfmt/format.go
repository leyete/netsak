// This file is part of Netsak
// Copyright (C) 2018 matterbeam
// Use of this source is governed by a GPLv3

package nsfmt

import (
	"io"
	"os"
)

// State represents the current state passed to the renderers.  It provides
// access to the io.Writer interface and allows to set the style to be used
// when writing the text.
type State interface {
	// SetStyle sets the style for the text to be displayed.
	SetStyle(s Style)

	// GetStyle returns the current style being used to display the text with.
	GetStyle() (s Style)

	// Write is the function to call to insert text in the place of the tag.
	// It writes at most len(p).  Returns the number of bytes written and
	// an error, if any.  The default style will be used.
	Write(p []byte) (n int, err error)

	// WriteString is like Write but allows to write a string directly.
	WriteString(s string) (n int, err error)
}

// Renderer is the interface that wraps the Render method.
//
// Render updates the supplied state according to the tag.
type Renderer interface {
	Render(st State, tag string)
}

func Fprintr(w io.Writer, r Renderer, s string) (n int, err error) {
	p := newPrinter()
	Parse(p, r, s)
	n, err = w.Write(p.Bytes())
	p.free()
	return
}

var defaultRenderer = DefaultRenderer{}

func Fprint(w io.Writer, s string) (n int, err error) {
	return Fprintr(w, defaultRenderer, s)
}

func Print(s string) (n int, err error) {
	return Fprint(os.Stdout, s)
}

// DefaultRenderer is the renderer that is installed by default on
// new views. It provides support for the following tags.
//
// {b}   - set bold text
// {/b}  - unset bold text
// {u}	 - set underlined text
// {/u}  - unset underlined text
// {d}   - set dim text
// {/d}  - unset dim text
// {bl}  - set blinking text
// {/bl} - unset blinking text
// {r}   - set reverse text (swap foreground and background colors)
// {/r}  - unset reverse text
//
// {black}   - set black foreground color
// {red}     - set red foreground color
// {green}   - set green foreground color
// {yellow}  - set yellow foreground color
// {blue}    - set blue foreground color
// {magenta} - set magenta foreground color
// {cyan}    - set cyan foreground color
// {white}   - set white foreground color
// {/fg}     - set default foreground color
// {/bg}     - set default background color
type DefaultRenderer struct{}

// Render renderizes the supplied tag accordingly.
func (r DefaultRenderer) Render(s State, tag string) {
	if f, ok := drTagMap[tag]; ok {
		f(s)
	}
}

// drTagMap is the default renderer's tag map.
var drTagMap = map[string]func(State){
	"b": func(s State) {
		s.SetStyle(s.GetStyle() | AttrBold)
	},
	"/b": func(s State) {
		s.SetStyle(s.GetStyle() &^ AttrBold)
	},
	"u": func(s State) {
		s.SetStyle(s.GetStyle() | AttrUnder)
	},
	"/u": func(s State) {
		s.SetStyle(s.GetStyle() &^ AttrUnder)
	},
	"d": func(s State) {
		s.SetStyle(s.GetStyle() | AttrDim)
	},
	"/d": func(s State) {
		s.SetStyle(s.GetStyle() &^ AttrDim)
	},
	"bl": func(s State) {
		s.SetStyle(s.GetStyle() | AttrBlink)
	},
	"/bl": func(s State) {
		s.SetStyle(s.GetStyle() &^ AttrBlink)
	},
	"r": func(s State) {
		s.SetStyle(s.GetStyle() | AttrReverse)
	},
	"/r": func(s State) {
		s.SetStyle(s.GetStyle() &^ AttrReverse)
	},
	"black": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorBlack)
	},
	"red": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorRed)
	},
	"green": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorGreen)
	},
	"yellow": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorYellow)
	},
	"blue": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorBlue)
	},
	"magenta": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorMagenta)
	},
	"cyan": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorCyan)
	},
	"white": func(s State) {
		s.SetStyle((s.GetStyle() & BitMaskFg) | ColorWhite)
	},
	"/fg": func(s State) {
		s.SetStyle(s.GetStyle() &^ FlagFg)
	},
	"/bg": func(s State) {
		s.SetStyle(s.GetStyle() &^ FlagBg)
	},
}
