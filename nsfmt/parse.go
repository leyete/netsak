package nsfmt

import (
	"bytes"
	"sync"
	"unicode/utf8"
)

// ----------------------------------------
// Default renderer state.
// ----------------------------------------

type pp struct {
	bytes.Buffer
	s Style
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
	// TODO: reset style (need an implementation of "deafult style")
	ppFree.Put(p)
}

func (p *pp) SetStyle(s Style) {
	p.WriteString(s.ANSI(ColorMode8Bit))
	p.s = s
}

func (p *pp) GetStyle() Style {
	return p.s
}

// ----------------------------------------

// Parse parses the tags contained in s.  The state st will be updated
// using r as the renderer.
func Parse(st State, r Renderer, s string) {
	end := len(s)
	lasti := 0 // Last index
	intag := false
	escape := false
	for i, c := range s {

		if escape {
			lasti = i
			escape = false
			continue
		}

		switch c {
		case '{':
			if !intag {
				intag = true
				if i > lasti {
					st.WriteString(s[lasti:i])
				}
				lasti = i + utf8.RuneLen(c) // Skip the '{'
				continue
			}
		case '}':
			if intag {
				intag = false
				if i > lasti {
					r.Render(st, s[lasti:i])
				}
				lasti = i + utf8.RuneLen(c) // Skip the '}'
				continue
			}
		case '\\':
			if i > lasti {
				st.WriteString(s[lasti:i])
			}
			// Escape the next rune.
			escape = true
		}

	}

	// Write the remaining text.  Skip any incomplete tag.
	if lasti < end && !intag {
		st.WriteString(s[lasti:])
	}
}

/*
func Parse(st State, r Renderer, s string) {
	tag := make([]byte, 10)[:0]
	intag := false
	escape := false

	for _, c := range s {

		if !escape {

			switch c {
			case '{':
				if !intag {
					intag = true
					continue
				}
			case '}':
				if intag {
					intag = false
					r.Render(st, string(tag))
					tag = tag[:0]
					continue
				}
			case '\\': // Escape the next rune
				escape = true
				continue
			}

		} else {
			escape = false
		}

		if intag {
			if c < utf8.RuneSelf {
				tag = append(tag, byte(c))
			} else {

			}
		} else {
			st.Write(RunesToUTF8(c))
		}

	}
}
*/
