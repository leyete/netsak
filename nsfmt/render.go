// This file is part of Netsak
// Copyright (C) 2018 matterbeam
// Use of this source is governed by a GPLv3

package nsfmt

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// ----------------------------------------
// Default renderer.
// ----------------------------------------

var renderer *BaseRenderer

// pp is the default renderer state, used in the functions provided by
// this package.  It is reused with sync.Pool to avoid allocations, this
// approach is the same as for the standard fmt package.
type pp struct {
	bytes.Buffer
	s SGR
}

var ppFree = sync.Pool{
	New: func() interface{} { return new(pp) },
}

func newPrinter() *pp {
	return ppFree.Get().(*pp)
}

func (p *pp) free() {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if p.Buffer.Cap() > 64<<10 {
		return
	}

	p.Buffer.Reset()
	p.s = DefaultSGR
	ppFree.Put(p)
}

func (p *pp) SetStyle(s SGR) {
	p.WriteString(s.Decode(ColorMode4Bit))
	p.s = s
}

func (p *pp) GetStyle() SGR {
	return p.s
}

// ----------------------------------------

// State represents the current state passed to the renderers.  It provides
// access to the io.Writer interface and allows to set the SGR to be used
// when displaying the text.
type State interface {
	// SetStyle sets the SGR for the text to be displayed.
	SetStyle(s SGR)

	// GetStyle returns the current SGR being used to display the text with.
	GetStyle() (s SGR)

	// Write is the function to call to update the current state.
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

// Fprintr renders s according to a renderer and writes to w.  It returns
// the number of bytes written and any encountered error.
func Fprintr(w io.Writer, r Renderer, s string) (n int, err error) {
	p := newPrinter()
	Parse(p, r, s)
	n, err = w.Write(p.Bytes())
	p.free()
	return
}

// Fprint renders s according to the default renderer and writes to w.
// It returns the number of bytes written and any encountered error.
func Fprint(w io.Writer, s string) (n int, err error) {
	return Fprintr(w, renderer, s)
}

// Print renders s according to the default renderer and writes to
// standard output.  It returns the number of bytes written and any
// encountered error.
func Print(s string) (n int, err error) {
	return Fprint(os.Stdout, s)
}

// Sprintr renders s according to a renderer and returns the resulting string.
func Sprintr(r Renderer, s string) string {
	p := newPrinter()
	Parse(p, r, s)
	str := p.String()
	p.free()
	return str
}

// Sprint renders s according to the default renderer anb returns
// the resulting string.
func Sprint(s string) string {
	return Sprintr(renderer, s)
}

// BaseRenderer is a basic value to implement renderers.  It provides access
// to concurrent-safe mechanisms to register/delete tags.
type BaseRenderer struct {
	mux sync.RWMutex
	m   map[string]func(State)
}

// Render renderizes the supplied tag accordingly.
func (r *BaseRenderer) Render(s State, tag string) {
	r.mux.RLock()
	if f, ok := r.m[tag]; ok {
		f(s)
	}
	r.mux.RUnlock()
}

// Register registers a function to be called when a tag needs to
// be rendered.
func (r *BaseRenderer) Register(tag string, f func(State)) {
	r.mux.Lock()
	r.m[tag] = f
	r.mux.Unlock()
}

// Delete removes the function associated with tag.  Future calls
// to Render won't be able to utilize this funciton to render the tag.
func (r *BaseRenderer) Delete(tag string) {
	r.mux.Lock()
	delete(r.m, tag)
	r.mux.Unlock()
}

func init() {
	// Initialize the default renderer with these tags.
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
	renderer = &BaseRenderer{
		m: map[string]func(State){
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
				s.SetStyle(s.GetStyle() | AttrDefFg)
			},
			"/bg": func(s State) {
				s.SetStyle(s.GetStyle() | AttrDefBg)
			},
		}}
}
