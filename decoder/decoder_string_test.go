package decoder

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

func encodedHtml() []byte {
	return []byte(`"\u003csection class=\"bg-gray solutions\" id=\"solutions\"\u003e\u003cdiv class=\"wrap text-center\"\u003e\u003ch2\u003eCUSTOM SOLUTIONS FOR YOUR BUSINESS\u003c/h2\u003e\u003cp\u003eAt ServerMania, we know the importance of choosing a server that is perfectly tailored for\u0026nbsp;\u003cbr class=\"hidden-lg hidden-xl\"\u003eyour needs. Every server application is unique and requires an expert understanding of\u0026nbsp;\u003cbr class=\"hidden-lg hidden-xl\"\u003ewhich components will deliver the best performance.\u003c/p\u003e\u003cdiv class=\"solutions-items-wrap\"\u003e\u003cdiv class=\"grid-wrapper\"\u003e\u003cpicture\u003e\u003c!--[--\u003e\u003c!--]--\u003e\u003cimg alt=\"portret Cam\" class=\"custom-image lazyloaded\" height=\"472\" width=\"308\" data-src=\"https://cdn.servermania.com/pictures/portret-cam-new.webp\" src=\"https://cdn.servermania.com/pictures/portret-cam-new.webp\"\u003e\u003c/picture\u003e\u003c/div\u003e\u003cdiv class=\"solutions-items text-left\"\u003e\u003cdiv\u003e\u003ch3\u003e\u003csvg class=\"c-icon\" height=\"40\" width=\"40\"\u003e\u003cuse href=\"#i-gamepad\"\u003e\u003c/use\u003e\u003c/svg\u003e Game Servers\u003c/h3\u003e\u003cp\u003eNever miss a second of the action with lightning\u003cbr class=\"hidden-sm hidden-md\"\u003efast connections with zero-lag server\u003cbr class=\"hidden-sm hidden-md\"\u003ecomponents.\u003c/p\u003e\u003ca href=\"/solutions/game-server-hosting.htm\" class=\"btn btn-sm btn-transparent\"\u003e\u003c!--[--\u003e LEARN MORE \u003csvg class=\"c-icon\" height=\"12\" width=\"12\"\u003e\u003cuse href=\"#i-arrow-right\"\u003e\u003c/use\u003e\u003c/svg\u003e \u003c!--]--\u003e\u003c/a\u003e\u003c/div\u003e\u003cdiv\u003e\u003ch3\u003e\u003csvg class=\"c-icon\" height=\"40\" width=\"40\"\u003e\u003cuse href=\"#i-hdd\"\u003e\u003c/use\u003e\u003c/svg\u003e Storage Servers\u003c/h3\u003e\u003cp\u003eEnterprise-grade components with easy to understand pricing, delivering an affordable storage solution for any project.\u003c/p\u003e\u003ca href=\"/solutions/storage-server-hosting.htm\" class=\"btn btn-sm btn-transparent\"\u003e\u003c!--[--\u003e LEARN MORE \u003csvg class=\"c-icon\" height=\"12\" width=\"12\"\u003e\u003cuse href=\"#i-arrow-right\"\u003e\u003c/use\u003e\u003c/svg\u003e \u003c!--]--\u003e\u003c/a\u003e\u003c/div\u003e\u003cdiv\u003e\u003ch3\u003e\u003csvg class=\"c-icon icon\" height=\"40\" width=\"40\"\u003e\u003cuse href=\"#i-server-2\"\u003e\u003c/use\u003e\u003c/svg\u003e VPN Server Solutions\u003c/h3\u003e\u003cp\u003eReduce IP and bandwidth costs with a high-performance server solution that empowers VPN providers to compete globally.\u003c/p\u003e\u003ca href=\"/solutions/vpn-server-hosting.htm\" class=\"btn btn-sm btn-transparent\"\u003e\u003c!--[--\u003e LEARN MORE \u003csvg class=\"c-icon\" height=\"12\" width=\"12\"\u003e\u003cuse href=\"#i-arrow-right\"\u003e\u003c/use\u003e\u003c/svg\u003e \u003c!--]--\u003e\u003c/a\u003e\u003c/div\u003e\u003cdiv\u003e\u003ch3\u003e\u003csvg class=\"c-icon icon\" height=\"40\" width=\"40\"\u003e\u003cuse href=\"#i-app-cart\"\u003e\u003c/use\u003e\u003c/svg\u003e E-Commerce Servers\u003c/h3\u003e\u003cp\u003eFeaturing a 100% network uptime guarantee and global connectivity to serve each customer to your business flawlessly.\u003c/p\u003e\u003cbr class=\"hidden-lg hidden-xl hidden-sm hidden-md\"\u003e\u003ca href=\"/solutions/ecommerce.htm\" class=\"btn btn-sm btn-transparent\"\u003e\u003c!--[--\u003e LEARN MORE \u003csvg class=\"c-icon\" height=\"12\" width=\"12\"\u003e\u003cuse href=\"#i-arrow-right\"\u003e\u003c/use\u003e\u003c/svg\u003e \u003c!--]--\u003e\u003c/a\u003e\u003c/div\u003e\u003c/div\u003e\u003c/div\u003e\u003c/div\u003e\u003c/section\u003e"`)
}

func encodedString(size int) []byte {
	data := "hello world"
	res := make([]byte, 0, size*len(data)+2)
	res = append(res, '"')
	for i := 0; i < size; i++ {
		res = append(res, data...)
	}
	res = append(res, '"')
	return res

}

func TestDecode_String_HTML(t *testing.T) {
	html := encodedHtml()

	EqualUnmarshaling[string](t, html)
}

func TestDecode_String(t *testing.T) {
	str := encodedString(1000)

	EqualUnmarshaling[string](t, str)
}

var benchStrHTML = encodedHtml()
var benchStr = encodedString(1000)

func BenchmarkString_HTML_Blaze(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		data := encodedHtml()
		err := Unmarshal(data, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStrHTML)))
}

func BenchmarkString_HTML_Std(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(benchStrHTML, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStrHTML)))
}

func BenchmarkString_HTML_GoJson(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(benchStrHTML, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStrHTML)))
}

func BenchmarkString_HTML_JsonIter(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(benchStrHTML, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStrHTML)))
}

func BenchmarkString_Simple_Blaze(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := Unmarshal(benchStr, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStr)))
}

func BenchmarkString_Simple_Std(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(benchStr, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStr)))
}

func BenchmarkString_Simple_GoJson(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(benchStr, &str)
		if err != nil {
			b.Fatal(err)
		}

	}
	b.SetBytes(int64(len(benchStr)))
}

func BenchmarkString_Simple_JsonIter(b *testing.B) {
	var str string
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(benchStr, &str)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(benchStr)))
}
