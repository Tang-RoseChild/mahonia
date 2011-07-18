include $(GOROOT)/src/Make.inc

TARG=mahonia

GOFILES=charset.go \
	utf8.go \
	utf16.go \
	ASCII.go \
	8bit.go \
	big5-data.go \
	big5.go \
	jis0201-data.go \
	jis0208-data.go \
	shiftjis.go \
	convert_string.go \
	reader.go \
	writer.go \
	jis0212-data.go \
	mbcs.go \
	euc-jp.go \
	translate.go \
	fallback.go \
	entity_data.go \
	entity.go \

include $(GOROOT)/src/Make.pkg