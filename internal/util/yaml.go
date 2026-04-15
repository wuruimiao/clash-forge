package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YamlOrderedMap 包装了 yaml.Node，用于构建有序的 YAML 映射
type YamlOrderedMap struct {
	Node *yaml.Node
}

// NewYamlOrderedMap 初始化一个 Mapping 类型的 YAML 节点
func NewYamlOrderedMap() *YamlOrderedMap {
	return &YamlOrderedMap{
		Node: &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		},
	}
}

// Add 向 Mapping 节点按顺序添加键值对
func (y *YamlOrderedMap) Add(key string, value any) error {
	// 1. 创建并添加 Key 节点
	y.Node.Content = append(y.Node.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
	})

	// 2. 创建 Value 节点并编码数据
	valNode := new(yaml.Node)
	if err := valNode.Encode(value); err != nil {
		return err
	}

	// 3. 将编码后的节点添加到 Content 中
	y.Node.Content = append(y.Node.Content, valNode)
	return nil
}

func UnmarshalYamlFile[T any](f string) (T, error) {
	var result T
	data, err := os.ReadFile(f)
	if err != nil {
		return result, err
	}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("parse region yaml: %w", err)
	}
	return result, nil
}
