package decoder

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/deveox/blaze/scopes"
)

func UnmarshalScoped(data []byte, v any, context scopes.Context, scope scopes.Decoding) error {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	t.ContextScope = context
	t.OperationScope = scope
	return t.Decode(v)
}

func Unmarshal(data []byte, v any) error {
	t := NewDecoder(data)
	defer decoderPool.Put(t)
	return t.Decode(v)
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

type Decoder struct {
	Buf            []byte
	ptr            unsafe.Pointer
	pos            int64
	start          int64
	ContextScope   scopes.Context
	OperationScope scopes.Decoding
}

func (t *Decoder) decode(v reflect.Value) error {

	fn, err := Decoders.Get(v)
	if err != nil {
		return err
	}
	return fn(t, v)
}
func (t *Decoder) nativeDecoder(v reflect.Value) error {
	c := t.char(t.ptr, t.pos)
	switch c {
	case '"':
		return t.decodeString(v)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
		t.Skip()
		return t.decodeNumber(v)
	case 't', 'f':
		t.Skip()
		return t.decodeBool(v)
	case 'n':
		t.Skip()
		return t.decodeNull(v)
	case '[':
		return t.decodeArrayOrSlice(v)
	case '{':
		return t.decodeObject(v)
	case TERMINATION_CHAR:
		return t.Error("[Blaze decode()] unexpected end of input, expected beginning of value")
	default:
		return t.Error("[Blaze decode()] invalid char, expected beginning of value")
	}
}

func (t *Decoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	return t.decode(rv)
}

func (t *Decoder) SkipWhitespace() {
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case ' ', '\t', '\n', '\r':
			t.pos++
		default:
			return
		}
	}
}

func (t *Decoder) Skip() error {
	for {
		c := t.char(t.ptr, t.pos)
		switch c {
		case ' ', '\t', '\n', '\r':
			t.pos++
		case '{':
			return t.SkipObject()
		case '[':
			return t.SkipArray()
		case '"':
			return t.SkipString()
		case 't':
			t.pos++
			return t.SkipTrue()
		case 'f':
			t.pos++
			return t.SkipFalse()
		case 'n':
			return t.SkipNull()
		case '0':
			t.start = t.pos
			return t.SkipZero()
		case '-':
			t.start = t.pos
			t.pos++
			c := t.char(t.ptr, t.pos)
			switch c {
			case '0':
				return t.SkipZero()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				t.pos++
				return t.SkipAnyNumber()
			default:
				return t.Error("[Blaze Skip()] invalid char, expected number")
			}

		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.start = t.pos
			t.pos++
			return t.SkipAnyNumber()
		case TERMINATION_CHAR:
			return t.Error("[Blaze Skip()] unexpected end of input, expected beginning of value")
		default:
			return t.Error("[Blaze Skip()] invalid char, expected beginning of value")
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
}

func (t *Decoder) Error(msg string) error {
	e := &Error{
		Message: msg,
		Offset:  int(t.pos),
	}
	areaStart := e.Offset - 20
	if areaStart < 0 {
		areaStart = 0
		e.AreaPos = e.Offset
	} else {
		e.AreaPos = 20
	}
	areaEnd := e.Offset + 20
	if areaEnd > len(t.Buf) {
		areaEnd = len(t.Buf)
	}
	e.Area = t.Buf[areaStart:areaEnd]
	return e
}

func (t *Decoder) ErrorF(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return t.Error(msg)
}

// char returns the byte at the given offset from the given pointer.
// it's similar to d.Buf[offset], but even though d.Buf[offset] never causes out-of-range access (because of the TERMINATION_CHAR),
// the compiler doesn't know that, so it can't optimize the bounds check away.
func (d *Decoder) char(ptr unsafe.Pointer, offset int64) byte {
	return *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(offset)))
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
