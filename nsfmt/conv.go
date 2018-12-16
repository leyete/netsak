


package nsfmt


// Itoa converts an integer to fixed-width decimal ASCII. If w is greater
// than zero it will serve as zero-padding size.
func Itoa(i, w int) []byte {
	// 20 is the maximum number of bytes needed to represent an int64
	// (from -9223372036854775808 to 9223372036854775807).
	bs := make([]byte, 20)
	pos := len(bs) - 1
	sig := i < 0

	// Truncate the width if necessary.
    if w > 20 {
        w = 20
    }

    // Adjust for negative numbers.
	if sig {
		i = -i
		w--    // The '-' takes one place in the padding
	}

    // Process the number from right to left by tens.
	for i >= 10 || w > 1 {
		q := i / 10
		bs[pos], i = byte('0'+i%10), q
		pos--
		w--
	}

    // Include any reminder.
	bs[pos] = byte('0' + i)

    // Include negative sign on negative numbers.
	if sig {
		pos--
		bs[pos] = '-'
	}

	return bs[pos:]
}
