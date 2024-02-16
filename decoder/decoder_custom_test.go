package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Settings struct {
	Name string
}

func (s *Settings) UnmarshalBlaze(d *Decoder, data []byte) error {
	*s = Settings{Name: "test"}
	return nil
}

func TestDecode_Custom(t *testing.T) {
	data := []byte(`null`)
	var s Settings
	err := DDecoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, "test", s.Name)
}
