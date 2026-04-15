package model

import (
	"fmt"
	"log/slog"
)

type ParsedConfig struct {
	Nodes []*Node // 解析出的代理节点列表
	*DNSRules
}

func (r *ParsedConfig) GetNodes() *Nodes {
	var (
		normals []*Node
		infos   []*Node
	)
	existNodeNameNo := make(map[string]int)
	existNodeUIDS := make(map[string]bool)

	for _, node := range r.Nodes {
		if node.IsInfo {
			infos = append(infos, node)
			continue
		}

		// 忽略重复的节点
		if existNodeUIDS[node.UID()] {
			slog.Warn("duplicated node, will ignore", "name", node.Name, "uid", node.UID())
			continue
		}
		existNodeUIDS[node.UID()] = true

		// 处理重名的节点
		if existNodeNameNo[node.Name] > 0 {
			node.Name = fmt.Sprintf("%s-d%d", node.Name, existNodeNameNo[node.Name])
		}
		existNodeNameNo[node.Name]++
		normals = append(normals, node)
	}
	return &Nodes{
		InfoNodes:   infos,
		NormalNodes: normals,
	}
}
