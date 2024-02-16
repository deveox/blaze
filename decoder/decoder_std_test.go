package decoder

import (
	"encoding/json"
	"testing"
	"time"
)

type A string

func (a *A) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*a = A(str)
	return nil
}

func TestUnmarshaler_Alias_Std(t *testing.T) {
	data := []byte(`"hello"`)
	EqualUnmarshaling[A](t, data)
}

type TTime struct {
	AddressableTime time.Time `json:"addressableTime,omitempty"`
}

func TestUnmarshaler_Time(t *testing.T) {
	data := []byte(`"2023-01-01T00:00:00Z"`)
	EqualUnmarshaling[time.Time](t, data)
	EqualUnmarshaling[**time.Time](t, data)
	data = []byte(`{"addressableTime":"2023-01-01T00:00:00Z"}`)
	EqualUnmarshaling[TTime](t, data)
}
