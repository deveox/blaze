package encoder

import (
	"testing"
	"time"
)

func TestEncode_Std(t *testing.T) {
	tm := time.Now()
	EqualMarshaling(t, tm)
}
