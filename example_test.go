package idgenx_test

import (
	"fmt"

	"github.com/lcylpzls/idgenx"
)

// ExampleGenerator 演示雪花 ID 生成与解析。
func ExampleGenerator() {
	g, err := idgenx.New(idgenx.DefaultConfig())
	if err != nil {
		panic(err)
	}
	id, err := g.Next()
	if err != nil {
		panic(err)
	}
	parts, err := g.Parse(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("节点 %d\n", parts.NodeID)
	// Output: 节点 0
}
