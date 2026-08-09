// Package main 演示 idgenx 基础用法：雪花 ID 与短 ID。
package main

import (
	"fmt"

	"github.com/lcylpzls/idgenx"
	"github.com/lcylpzls/idgenx/shortid"
)

// run 演示雪花与短 ID 生成。
func run() error {
	g, err := idgenx.New(idgenx.DefaultConfig())
	if err != nil {
		return err
	}
	id, err := g.Next()
	if err != nil {
		return err
	}
	parts, err := g.Parse(id)
	if err != nil {
		return err
	}
	fmt.Printf("雪花 ID：%d（节点 %d，序列 %d）\n", id, parts.NodeID, parts.Sequence)

	code, err := shortid.Generate(8)
	if err != nil {
		return err
	}
	fmt.Printf("短 ID：%s\n", code)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("示例失败：", err)
	}
}
