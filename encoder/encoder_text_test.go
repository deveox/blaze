package encoder

import "testing"

type TEncoder struct{}

func (t TEncoder) MarshalText() ([]byte, error) {
	return []byte("test"), nil
}

func TestEncode_Text(t *testing.T) {
	te := TEncoder{}
	EqualMarshaling(t, te)
}
