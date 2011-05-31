package mahonia

// Converters for the Shift-JIS encoding.

import (
	"sync"
)

func init() {
	RegisterCharset(&Charset{
		Name:    "Shift_JIS",
		Aliases: []string{"MS_Kanji", "csShiftJIS", "SJIS"},
		NewDecoder: func() Decoder {
			sjisOnce.Do(makeSjisTables)
			return decodeSjisRune
		},
		NewEncoder: func() Encoder {
			sjisOnce.Do(makeSjisTables)
			return encodeSjisRune
		},
	})
}

func decodeSjisRune(p []byte) (rune, size int, status Status) {
	if len(p) == 0 {
		status = NO_ROOM
		return
	}

	b := p[0]

	rune = int(sjisToUnicode[b])
	if rune != 0 || b == 0 {
		return rune, 1, SUCCESS
	}

	if len(p) < 2 {
		status = NO_ROOM
		return
	}

	rune = int(sjisToUnicode[int(b)<<8+int(p[1])])
	if rune != 0 {
		return rune, 2, SUCCESS
	}

	return 0xfffd, 1, INVALID_CHAR
}

func encodeSjisRune(p []byte, rune int) (size int, status Status) {
	if len(p) == 0 {
		status = NO_ROOM
		return
	}

	if rune > 0xffff {
		p[0] = '?'
		return 1, INVALID_CHAR
	}

	c := unicodeToSjis[rune]
	if c == 0 && rune != 0 {
		p[0] = '?'
		return 1, INVALID_CHAR
	}

	if c < 256 {
		p[0] = byte(c)
		return 1, SUCCESS
	}

	if len(p) < 2 {
		status = NO_ROOM
		return
	}

	p[0] = byte(c >> 8)
	p[1] = byte(c)
	return 2, SUCCESS
}

var sjisOnce sync.Once

var sjisToUnicode []uint16
var unicodeToSjis []uint16

func makeSjisTables() {
	sjisToUnicode = make([]uint16, 65536)
	unicodeToSjis = make([]uint16, 65536)

	for jis0208, unicode := range jis0208ToUnicode {
		if unicode == 0 {
			continue
		}

		j1 := jis0208 >> 8
		j2 := jis0208 & 0xff

		var s1, s2 int // the bytes of the shift-jis code

		if j1 < 95 {
			s1 = (j1+1)/2 + 112
		} else {
			s1 = (j1+1)/2 + 176
		}

		if j1&1 == 1 {
			s2 = j2 + 31
			if j2 >= 96 {
				s2++
			}
		} else {
			s2 = j2 + 126
		}

		sjis := s1<<8 + s2

		sjisToUnicode[sjis] = uint16(unicode)
		unicodeToSjis[unicode] = uint16(sjis)
	}

	for jis0201, unicode := range jis0201ToUnicode {
		if unicode == 0 {
			continue
		}

		sjisToUnicode[jis0201] = uint16(unicode)
		unicodeToSjis[unicode] = uint16(jis0201)
	}

	for i := 0; i < 32; i++ {
		sjisToUnicode[i] = uint16(i)
	}

	sjisToUnicode[127] = 127
}
