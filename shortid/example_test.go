package shortid_test

import (
	"fmt"

	"github.com/lcylpzls/idgenx/shortid"
)

// ExampleGenerate 演示短 ID 生成。
func ExampleGenerate() {
	code, err := shortid.Generate(8)
	if err != nil {
		panic(err)
	}
	fmt.Printf("长度 %d 字符集合法 %v\n", len(code), shortid.IsValid(code, shortid.AlphabetBase62))
	// Output: 长度 8 字符集合法 true
}
