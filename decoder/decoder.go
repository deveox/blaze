package decoder

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/deveox/blaze/ctx"
	"github.com/deveox/blaze/scopes"
)

// In order to decode, you have to traverse the input buffer character by position. At that time, if you check whether the buffer has reached the end, it will be very slow.
// Therefore, by adding the NUL (\000) character to the end of the read buffer as shown below, it is possible to check the termination character at the same time as other characters.
const TERMINATION_CHAR = '\000'
const MAX_DEPTH = 1000

type Decoder struct {
	*ctx.Ctx
	config        *Config
	depth         int
	Buf           []byte
	ptr           unsafe.Pointer
	pos           int64
	start         int64
	operation     scopes.Decoding
	Changes       []string
	ChangesPrefix string
}

func (d *Decoder) Unmarshal(data []byte, v any) error {
	return d.config.UnmarshalScopedCtx(data, v, d.operation, d.Ctx)
}

func (d *Decoder) UnmarshalScoped(data []byte, v any, operation scopes.Decoding) error {
	return d.config.UnmarshalScopedCtx(data, v, operation, d.Ctx)
}

func (d *Decoder) UnmarshalScopedWithChanges(data []byte, v any, operation scopes.Decoding) ([]string, error) {
	return d.config.UnmarshalScopedWithChangesCtx(data, v, operation, d.Ctx)
}

func (d *Decoder) UnmarshalWithChanges(data []byte, v any) ([]string, error) {
	return d.config.UnmarshalScopedWithChangesCtx(data, v, d.operation, d.Ctx)
}

func (d *Decoder) Context() scopes.Context {
	return d.config.Scope
}

func (d *Decoder) Operation() scopes.Decoding {
	return d.operation
}

func (d *Decoder) decode(v reflect.Value) error {
	return getDecoderFn(v.Type())(d, v)

}

func (d *Decoder) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return d.Error("[Blaze decode()] can't decode to nil value")
	}
	if rv.Kind() != reflect.Pointer {
		return d.ErrorF("[Blaze decode()] can't decode to non-pointer value '%s'", rv.Type())
	}
	return d.decode(rv)
}

func (d *Decoder) SkipWhitespace() {
	for {
		c := d.char()
		switch c {
		case ' ', '\t', '\n', '\r':
			d.pos++
		default:
			return
		}
	}
}

func (d *Decoder) Skip() error {
	d.SkipWhitespace()
	var err error
	start := d.pos
	c := d.char()
	switch c {
	case '{':
		if d.depth > MAX_DEPTH {
			return d.Error("[Blaze decode()] maximum depth reached")
		}
		err = d.SkipObject()
		d.depth--
	case '[':
		if d.depth > MAX_DEPTH {
			return d.Error("[Blaze decode()] maximum depth reached")
		}
		err = d.SkipArray()
		d.depth--
	case '"':
		err = d.SkipString()
	case 't':
		err = d.SkipTrue()
	case 'f':
		err = d.SkipFalse()
	case 'n':
		err = d.ScanNull()
	case '0':
		err = d.SkipZero(true)
	case '-':
		err = d.SkipMinus(true)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err = d.SkipNumber(true, true)
	case TERMINATION_CHAR:
		return d.Error("[Blaze Skip()] unexpected end of input, expected beginning of value")
	default:
		return d.Error("[Blaze Skip()] invalid char, expected beginning of value")
	}
	d.start = start
	return err
}

func (d *Decoder) init(data []byte) {
	// Reset the buffer.
	d.Buf = d.Buf[:0]
	d.Buf = append(d.Buf, data...)
	// Add a termination character to the end of the buffer to avoid out-of-range access.
	d.Buf = append(d.Buf, TERMINATION_CHAR)
	// Set the pointer to the first byte of the buffer.
	d.ptr = unsafe.Pointer(unsafe.SliceData(d.Buf))
	d.pos = 0
	d.start = 0
	d.operation = 0
	d.ChangesPrefix = ""
	d.Changes = d.Changes[:0]
	d.depth = 0
}

func (d *Decoder) Error(msg string) error {
	e := &Error{
		Message: msg,
		Offset:  int(d.pos),
	}
	areaStart := e.Offset - 20
	if areaStart < 0 {
		areaStart = 0
		e.AreaPos = e.Offset
	} else {
		e.AreaPos = 20
	}
	areaEnd := e.Offset + 20
	if areaEnd > len(d.Buf) {
		areaEnd = len(d.Buf)
	}
	e.Area = d.Buf[areaStart:areaEnd]
	return e
}

func (d *Decoder) ErrorF(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return d.Error(msg)
}

// char returns the byte at the given offset from the given pointer.
// it's similar to d.Buf[offset], but even though d.Buf[offset] never causes out-of-range access (because of the TERMINATION_CHAR),
// the compiler doesn't know that, so it can't optimize the bounds check away.
func (d *Decoder) char() byte {
	return *(*byte)(unsafe.Pointer(uintptr(d.ptr) + uintptr(d.pos)))
}

func BytesToString(b []byte) string {
	// Ignore if your IDE shows an error here; it's a false positive.
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

type Error struct {
	Message string
	Offset  int
	Area    []byte
	AreaPos int
}

func (e *Error) Error() string {
	str := fmt.Sprintf("%s at position %d:\n%s\n", e.Message, e.Offset, string(e.Area))
	for i := 0; i < len(e.Area); i++ {
		if i == e.AreaPos {
			str += "^"
		} else {
			str += " "
		}
	}
	return str
}
