package model

import (
	"sort"

	"github.com/wuruimiao/clash-forge/internal/util"
)

type NodeGroupType int

const (
	Region NodeGroupType = iota
	Mult
)

// NodeGroup 表示一个自动生成的节点分组。
type NodeGroup struct {
	Name  string // 显示组名称：如 "香港"、"日本"、"倍率0.1"
	Type  NodeGroupType
	Nodes []*Node // 组内节点列表
}

type NodeGroups struct {
	Groups map[NodeGroupType][]*NodeGroup
}

func (r *NodeGroups) GroupNodeNames() *util.OrderedMap[string, []string] {
	result := util.NewOrderedMap[string, []string]()

	for _, gs := range r.Groups {
		for _, group := range gs {
			nodes, _ := result.Get(group.Name)
			for _, node := range group.Nodes {
				nodes = append(nodes, node.Name)
			}
			result.Set(group.Name, nodes)
		}
	}

	return result
}

func (r *NodeGroups) GroupMultNames() []string {
	var result []string
	for _, g := range r.Groups[Mult] {
		result = append(result, g.Name)
	}
	sort.Strings(result)
	return result
}
