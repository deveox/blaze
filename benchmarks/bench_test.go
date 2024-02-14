package benchmarks

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deveox/blaze/encoder"

	stdjson "encoding/json"

	"github.com/francoispqt/gojay"
	"github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"github.com/pquerna/ffjson/ffjson"
	segmentiojson "github.com/segmentio/encoding/json"
	"github.com/wI2L/jettison"
)

type codeResponse struct {
	Tree     *codeNode `json:"tree"`
	Username string    `json:"username"`
}

type codeNode struct {
	Name     string      `json:"name"`
	Kids     []*codeNode `json:"kids"`
	CLWeight float64     `json:"cl_weight"`
	Touches  int         `json:"touches"`
	MinT     int64       `json:"min_t"`
	MaxT     int64       `json:"max_t"`
	MeanT    int64       `json:"mean_t"`
}

var codeJSON []byte
var codeStruct codeResponse

func codeInit() {
	f, err := os.Open("./code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	if err := stdjson.Unmarshal(codeJSON, &codeStruct); err != nil {
		panic("unmarshal code.json: " + err.Error())
	}
	{
		stdjsonbytes, err := stdjson.Marshal(&codeStruct)
		if err != nil {
			panic("marshal code.json: " + err.Error())
		}
		jsonbytes, err := json.Marshal(&codeStruct)
		if err != nil {
			panic("marshal code.json: " + err.Error())
		}
		if len(stdjsonbytes) != len(jsonbytes) {
			panic(fmt.Sprintf("stdjson = %d but go-json = %d", len(stdjsonbytes), len(jsonbytes)))
		}
	}
	if _, err := json.Marshal(&codeStruct); err != nil {
		panic("marshal code.json: " + err.Error())
	}
	if !bytes.Equal(data, codeJSON) {
		println("different lengths", len(data), len(codeJSON))
		for i := 0; i < len(data) && i < len(codeJSON); i++ {
			if data[i] != codeJSON[i] {
				println("re-marshal: changed at byte", i)
				println("orig: ", string(codeJSON[i-10:i+10]))
				println("new: ", string(data[i-10:i+10]))
				break
			}
		}
		panic("re-marshal code.json: different result")
	}
}

func Benchmark_MarshalBigData_GoJson(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := json.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func Benchmark_MarshalBigData_EncodingJson(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := stdjson.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func Benchmark_MarshalBigData_JsonIter(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := json.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func Benchmark_MarshalBigData_Jettison(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := jettison.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func Benchmark_MarshalBigData_SegmentioJson(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := segmentiojson.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func Benchmark_MarshalBigData_Blaze(b *testing.B) {
	b.ReportAllocs()
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Marshal(&codeStruct); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	})
	b.SetBytes(int64(len(codeJSON)))
}

func benchMarshalBytes(n int, marshaler func(interface{}) ([]byte, error)) func(*testing.B) {
	sample := []byte("hello world")
	// Use a struct pointer, to avoid an allocation when passing it as an
	// interface parameter to Marshal.
	v := &struct {
		Bytes []byte
	}{
		bytes.Repeat(sample, (n/len(sample))+1)[:n],
	}
	return func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := marshaler(v); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
	}
}

func Benchmark_MarshalBytes_EncodingJson(b *testing.B) {
	// 32 fits within encodeState.scratch.
	b.Run("32", benchMarshalBytes(32, stdjson.Marshal))
	// 256 doesn't fit in encodeState.scratch, but is small enough to
	// allocate and avoid the slower base64.NewEncoder.
	b.Run("256", benchMarshalBytes(256, stdjson.Marshal))
	// 4096 is large enough that we want to avoid allocating for it.
	b.Run("4096", benchMarshalBytes(4096, stdjson.Marshal))
}

func Benchmark_MarshalBytes_JsonIter(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// 32 fits within encodeState.scratch.
	b.Run("32", benchMarshalBytes(32, json.Marshal))
	// 256 doesn't fit in encodeState.scratch, but is small enough to
	// allocate and avoid the slower base64.NewEncoder.
	b.Run("256", benchMarshalBytes(256, json.Marshal))
	// 4096 is large enough that we want to avoid allocating for it.
	b.Run("4096", benchMarshalBytes(4096, json.Marshal))
}

func Benchmark_MarshalBytes_GoJson(b *testing.B) {
	b.ReportAllocs()
	// 32 fits within encodeState.scratch.
	b.Run("32", benchMarshalBytes(32, json.Marshal))
	// 256 doesn't fit in encodeState.scratch, but is small enough to
	// allocate and avoid the slower base64.NewEncoder.
	b.Run("256", benchMarshalBytes(256, json.Marshal))
	// 4096 is large enough that we want to avoid allocating for it.
	b.Run("4096", benchMarshalBytes(4096, json.Marshal))
}

func Benchmark_MarshalBytes_Blaze(b *testing.B) {
	b.ReportAllocs()
	// 32 fits within encodeState.scratch.
	b.Run("32", benchMarshalBytes(32, encoder.Marshal))
	// 256 doesn't fit in encodeState.scratch, but is small enough to
	// allocate and avoid the slower base64.NewEncoder.
	b.Run("256", benchMarshalBytes(256, encoder.Marshal))
	// 4096 is large enough that we want to avoid allocating for it.
	b.Run("4096", benchMarshalBytes(4096, encoder.Marshal))
}

func Benchmark_MarshalString_EncodingJson(b *testing.B) {
	b.ReportAllocs()
	j := struct {
		Bar string `json:"bar,string"`
	}{
		Bar: `foobar`,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := stdjson.Marshal(&j); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Benchmark_MarshalString_JsonIter(b *testing.B) {
	b.ReportAllocs()
	j := struct {
		Bar string `json:"bar,string"`
	}{
		Bar: `foobar`,
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := json.Marshal(&j); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Benchmark_MarshalString_GoJson(b *testing.B) {
	b.ReportAllocs()
	j := struct {
		Bar string `json:"bar,string"`
	}{
		Bar: `foobar`,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := json.Marshal(&j); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Benchmark_MarshalString_Blaze(b *testing.B) {
	b.ReportAllocs()
	j := struct {
		Bar string `json:"bar,string"`
	}{
		Bar: `foobar`,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Marshal(&j); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func Benchmark_Encode_SmallStruct_EncodingJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_FFJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := ffjson.Marshal(NewSmallPayloadFFJson()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_JsonIter(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_Jettison(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := jettison.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_GoJay(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := gojay.MarshalJSONObject(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_SegmentioJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := segmentiojson.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_GoJson(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Encode_SmallStruct_Blaze(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := encoder.Marshal(NewSmallPayload()); err != nil {
			b.Fatal(err)
		}
	}
}
