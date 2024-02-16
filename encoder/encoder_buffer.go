package encoder

func (e *Encoder) Reset() {
	e.bytes = e.bytes[:0]
}

func (e *Encoder) Write(b []byte) {
	e.bytes = append(e.bytes, b...)
}

func (e *Encoder) WriteString(s string) {
	e.bytes = append(e.bytes, s...)
}

func (e *Encoder) WriteByte(b byte) {
	e.bytes = append(e.bytes, b)
}
