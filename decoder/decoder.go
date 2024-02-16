package decoder

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/deveox/blaze/scopes"
)

func UnmarshalScoped(data []byte, v any, scope scopes.Decoding) error {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	t.operationScope = scope
	return t.Decode(v)
}

func UnmarshalScopedWithChanges(data []byte, v any, scope scopes.Decoding) ([]string, error) {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	t.operationScope = scope
	t.Changes = make([]string, 0, 10)
	err := t.Decode(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err
}

func Unmarshal(data []byte, v any) error {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	return t.Decode(v)
}

func UnmarshalWithChanges(data []byte, v any) ([]string, error) {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	t.Changes = make([]string, 0, 10)
	err := t.Decode(v)
	changes := t.Changes
	t.Changes = nil
	return changes, err

}

func NewDecoder(data []byte) *Decoder {
	if v := decoderPool.Get(); v != nil {
		d := v.(*Decoder)
		d.init(data)
		return d
	}
	t := &Decoder{}
	t.init(data)
	return t
}

// In order to decode, you have to traverse the input buffer character by position. At that time, if you check whether the buffer has reached the end, it will be very slow.
// Therefore, by adding the NUL (\000) character to the end of the read buffer as shown below, it is possible to check the termination character at the same time as other characters.
const TERMINATION_CHAR = '\000'
const MAX_DEPTH = 1000

type Decoder struct {
	depth          int
	Buf            []byte
	ptr            unsafe.Pointer
	pos            int64
	start          int64
	contextScope   scopes.Context
	operationScope scopes.Decoding
	Changes        []string
	ChangesPrefix  string
}

func (d *Decoder) Context() scopes.Context {
	return d.contextScope
}

func (d *Decoder) Operation() scopes.Decoding {
	return d.operationScope
}

func (d *Decoder) decode(v reflect.Value) error {
	return getDecoderFn(v.Type())(d, v)
}

func (d *Decoder) Decode(v any) error {
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
	for {
		c := d.char()
		switch c {
		case ' ', '\t', '\n', '\r':
			d.pos++
		case '{':
			if d.depth > MAX_DEPTH {
				return d.Error("[Blaze decode()] maximum depth reached")
			}
			err := d.SkipObject()
			d.depth--
			return err
		case '[':
			if d.depth > MAX_DEPTH {
				return d.Error("[Blaze decode()] maximum depth reached")
			}
			err := d.SkipArray()
			d.depth--
			return err
		case '"':
			return d.SkipString()
		case 't':
			d.pos++
			return d.SkipTrue()
		case 'f':
			d.pos++
			return d.SkipFalse()
		case 'n':
			return d.ScanNull()
		case '0':
			d.start = d.pos
			return d.SkipZero(true)
		case '-':
			d.start = d.pos
			return d.SkipMinus(true)
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d.start = d.pos
			return d.SkipNumber(true, true)
		case TERMINATION_CHAR:
			return d.Error("[Blaze Skip()] unexpected end of input, expected beginning of value")
		default:
			return d.Error("[Blaze Skip()] invalid char, expected beginning of value")
		}
	}
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
	d.operationScope = 0
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

// tokenPool has a pool of Decoder instances.
// It's used to reduce the number of allocations when decoding JSON.
var decoderPool sync.Pool

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
