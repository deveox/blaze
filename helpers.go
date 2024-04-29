package blaze

import (
	"fmt"
	"strconv"

	"github.com/deveox/blaze/decoder"
	"github.com/deveox/blaze/encoder"
)

type StringFloat float64

func (s StringFloat) MarshalBlaze(e *encoder.Encoder) error {
	return e.Encode(strconv.FormatFloat(float64(s), 'f', -1, 64))
}

func (s *StringFloat) UnmarshalBlaze(d *decoder.Decoder, data []byte) error {
	if data[0] == '"' {
		var str string
		err := d.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		if str == "" {
			*s = 0
			return nil
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*s = StringFloat(f)
		return nil
	}
	var f float64
	err := d.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	*s = StringFloat(f)
	return nil
}

type StringBool bool

func (b StringBool) MarshalBlaze(e *encoder.Encoder) error {
	if b {
		return e.Encode("true")
	}
	return e.Encode("false")
}

func (b *StringBool) UnmarshalBlaze(d *decoder.Decoder, data []byte) error {
	if data[0] == '"' {
		var str string
		err := d.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		switch str {
		case "true":
			*b = true
		case "false":
			*b = false
		default:
			return fmt.Errorf("StringBool can't unmarshal %s. Expected true or false", str)
		}
	}
	var bol bool
	err := d.Unmarshal(data, &bol)
	if err != nil {
		return err
	}
	*b = StringBool(bol)
	return nil
}
