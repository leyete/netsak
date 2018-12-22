package nsfmt

import (
	"unicode/utf8"
)

// Parse parses the tags contained in s.  The state st will be updated
// using r as the renderer.
func Parse(st State, r Renderer, s string) {
	end := len(s)
	lasti := 0
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
