// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yaml

const (
	// The size of the input raw buffer.
	inputRawBufferSize = 512

	// The size of the input buffer.
	// It should be possible to decode the whole raw buffer.
	inputBufferSize = inputRawBufferSize * 3

	// The size of the output buffer.
	outputBufferSize = 128

	// The size of the output raw buffer.
	// It should be possible to encode the whole output buffer.
	outputRawBufferSize = (outputBufferSize*2 + 2)

	// The size of other stacks and queues.
	initialStackSize  = 16
	initialQueueSize  = 16
	initialStringSize = 16
)

// Check if the character at the specified position is an alphabetical
// character, a digit, '_', or '-'.
func isAlpha(b []byte, i int) bool {
	return b[i] >= '0' && b[i] <= '9' || b[i] >= 'A' && b[i] <= 'Z' || b[i] >= 'a' && b[i] <= 'z' || b[i] == '_' || b[i] == '-'
}

// Check if the character at the specified position is a digit.
func isDigit(b []byte, i int) bool {
	return b[i] >= '0' && b[i] <= '9'
}

// Get the value of a digit.
func asDigit(b []byte, i int) int {
	return int(b[i]) - '0'
}

// Check if the character at the specified position is a hex-digit.
func isHex(b []byte, i int) bool {
	return b[i] >= '0' && b[i] <= '9' || b[i] >= 'A' && b[i] <= 'F' || b[i] >= 'a' && b[i] <= 'f'
}

// Get the value of a hex-digit.
func asHex(b []byte, i int) int {
	bi := b[i]
	if bi >= 'A' && bi <= 'F' {
		return int(bi) - 'A' + 10
	}
	if bi >= 'a' && bi <= 'f' {
		return int(bi) - 'a' + 10
	}
	return int(bi) - '0'
}

// Check if the character is ASCII.
func isASCII(b []byte, i int) bool {
	return b[i] <= 0x7F
}

// Check if the character at the start of the buffer can be printed unescaped.
func isPrintable(b []byte, i int) bool {
	return ((b[i] == 0x0A) || // . == #x0A
		(b[i] >= 0x20 && b[i] <= 0x7E) || // #x20 <= . <= #x7E
		(b[i] == 0xC2 && b[i+1] >= 0xA0) || // #0xA0 <= . <= #xD7FF
		(b[i] > 0xC2 && b[i] < 0xED) ||
		(b[i] == 0xED && b[i+1] < 0xA0) ||
		(b[i] == 0xEE) ||
		(b[i] == 0xEF && // #xE000 <= . <= #xFFFD
			!(b[i+1] == 0xBB && b[i+2] == 0xBF) && // && . != #xFEFF
			!(b[i+1] == 0xBF && (b[i+2] == 0xBE || b[i+2] == 0xBF))))
}

// Check if the character at the specified position is NUL.
func isZ(b []byte, i int) bool {
	return b[i] == 0x00
}

// Check if the beginning of the buffer is a BOM.
func isBom(b []byte, i int) bool {
	return b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF
}

// Check if the character at the specified position is space.
func isSpace(b []byte, i int) bool {
	return b[i] == ' '
}

// Check if the character at the specified position is colon.
func isColon(b []byte, i int) bool {
	return b[i] == ':'
}

// Check if the character at the specified position is tab.
func isTab(b []byte, i int) bool {
	return b[i] == '\t'
}

// Check if the character at the specified position is blank (space or tab).
func isBlank(b []byte, i int) bool {
	//return is_space(b, i) || is_tab(b, i)
	return b[i] == ' ' || b[i] == '\t'
}

// Check if the character at the specified position is a line break.
func isBreak(b []byte, i int) bool {
	return (b[i] == '\r' || // CR (#xD)
		b[i] == '\n' || // LF (#xA)
		b[i] == 0xC2 && b[i+1] == 0x85 || // NEL (#x85)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA8 || // LS (#x2028)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA9) // PS (#x2029)
}

func isCrlf(b []byte, i int) bool {
	return b[i] == '\r' && b[i+1] == '\n'
}

// Check if the character is a line break or NUL.
func isBreakz(b []byte, i int) bool {
	//return is_break(b, i) || is_z(b, i)
	return (        // is_break:
	b[i] == '\r' || // CR (#xD)
		b[i] == '\n' || // LF (#xA)
		b[i] == 0xC2 && b[i+1] == 0x85 || // NEL (#x85)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA8 || // LS (#x2028)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA9 || // PS (#x2029)
		// is_z:
		b[i] == 0)
}

// Check if the character is a line break, space, or NUL.
func isSpacez(b []byte, i int) bool {
	//return is_space(b, i) || is_breakz(b, i)
	return ( // is_space:
	b[i] == ' ' ||
		// is_breakz:
		b[i] == '\r' || // CR (#xD)
		b[i] == '\n' || // LF (#xA)
		b[i] == 0xC2 && b[i+1] == 0x85 || // NEL (#x85)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA8 || // LS (#x2028)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA9 || // PS (#x2029)
		b[i] == 0)
}

// Check if the character is a line break, space, tab, or NUL.
func isBlankz(b []byte, i int) bool {
	//return is_blank(b, i) || is_breakz(b, i)
	return ( // is_blank:
	b[i] == ' ' || b[i] == '\t' ||
		// is_breakz:
		b[i] == '\r' || // CR (#xD)
		b[i] == '\n' || // LF (#xA)
		b[i] == 0xC2 && b[i+1] == 0x85 || // NEL (#x85)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA8 || // LS (#x2028)
		b[i] == 0xE2 && b[i+1] == 0x80 && b[i+2] == 0xA9 || // PS (#x2029)
		b[i] == 0)
}

// Determine the width of the character.
func width(b byte) int {
	// Don't replace these by a switch without first
	// confirming that it is being inlined.
	if b&0x80 == 0x00 {
		return 1
	}
	if b&0xE0 == 0xC0 {
		return 2
	}
	if b&0xF0 == 0xE0 {
		return 3
	}
	if b&0xF8 == 0xF0 {
		return 4
	}
	return 0

}
