package encoder

import "reflect"

func newCondAddrEncoder(canAddrEnc, elseEnc EncoderFn) EncoderFn {
	return func(e *Encoder, v reflect.Value) error {
		if v.CanAddr() {
			return canAddrEnc(e, v)
		}
		return elseEnc(e, v)
	}
}
