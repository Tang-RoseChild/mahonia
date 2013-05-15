package mahonia

// Converters for the Shift-JIS encoding.

import (
	"unicode/utf8"
)

func init() {
	RegisterCharset(&Charset{
		Name:    "Shift_JIS",
		Aliases: []string{"MS_Kanji", "csShiftJIS", "SJIS"},
		NewDecoder: func() Decoder {
			return decodeSJIS
		},
		NewEncoder: func() Encoder {
			jis0208Once.Do(reverseJIS0208Table)
			return encodeSJIS
		},
	})
}

func decodeSJIS(p []byte) (c rune, size int, status Status) {
	if len(p) == 0 {
		return 0, 0, NO_ROOM
	}

	b := p[0]
	if b == 0x7e {
		return '‾', 1, SUCCESS
	}
	if b == 0x5c {
		return '¥', 1, SUCCESS
	}
	if b < 0x80 {
		return rune(b), 1, SUCCESS
	}

	if 0xa1 <= b && b <= 0xdf {
		return rune(b) + (0xff61 - 0xa1), 1, SUCCESS
	}

	if b == 0x80 || b == 0xa0 || b >= 0xf0 {
		return utf8.RuneError, 1, INVALID_CHAR
	}

	if len(p) < 2 {
		return 0, 0, NO_ROOM
	}

	s1 := b
	s2 := p[1]

	var j1, j2 byte
	if s1 < 0xa0 {
		j1 = (s1 - 112) * 2
	} else {
		j1 = (s1 - 176) * 2
	}

	if s2 >= 0x9f {
		j2 = s2 - 126
	} else {
		j1--
		j2 = s2 - 31
		if s2 > 0x7f {
			j2--
		}
	}

	jis0208 := int(j1)<<8 + int(j2)
	unicode := jis0208ToUnicode[jis0208]
	if unicode == 0 {
		return utf8.RuneError, 2, INVALID_CHAR
	}
	return rune(unicode), 2, SUCCESS
}

func encodeSJIS(p []byte, c rune) (size int, status Status) {
	if len(p) == 0 {
		return 0, NO_ROOM
	}

	if c < 0x80 && c != '\\' && c != '~' {
		p[0] = byte(c)
		return 1, SUCCESS
	}

	if c == '‾' {
		p[0] = 0x7e
		return 1, SUCCESS
	}

	if c == '¥' {
		p[0] = 0x5c
		return 1, SUCCESS
	}

	if 0xff61 <= c && c <= 0xff9f {
		// half-width katakana
		p[0] = byte(c - (0xff61 - 0xa1))
		return 1, SUCCESS
	}

	if len(p) < 2 {
		return 0, NO_ROOM
	}

	if c > 0xffff {
		p[0] = '?'
		return 1, INVALID_CHAR
	}

	jis0208 := unicodeToJIS0208[c]
	if jis0208 == 0 {
		p[0] = '?'
		return 1, INVALID_CHAR
	}

	j1 := byte(jis0208 >> 8)
	j2 := byte(jis0208)

	if j1 < 95 {
		p[0] = (j1+1)/2 + 112
	} else {
		p[0] = (j1+1)/2 + 176
	}

	if j1&1 == 1 {
		p[1] = j2 + 31
		if j2 >= 96 {
			p[1]++
		}
	} else {
		p[1] = j2 + 126
	}

	return 2, SUCCESS
}
