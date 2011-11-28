package mahonia

import "unicode/utf8"

func init() {
	RegisterCharset(&Charset{
		Name:       "UTF-8",
		NewDecoder: func() Decoder { return decodeUTF8Rune },
		NewEncoder: func() Encoder { return encodeUTF8Rune },
	})
}

func decodeUTF8Rune(p []byte) (rune, size int, status Status) {
	if len(p) == 0 {
		status = NO_ROOM
		return
	}

	if p[0] < 128 {
		return int(p[0]), 1, SUCCESS
	}

	rune, size = utf8.DecodeRune(p)

	if rune == 0xfffd {
		if utf8.FullRune(p) {
			status = INVALID_CHAR
			return
		}

		return 0, 0, NO_ROOM
	}

	status = SUCCESS
	return
}

func encodeUTF8Rune(p []byte, rune int) (size int, status Status) {
	size = utf8.RuneLen(rune)
	if size > len(p) {
		return 0, NO_ROOM
	}

	return utf8.EncodeRune(p, rune), SUCCESS
}
