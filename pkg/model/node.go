package model

import (
	"fmt"
	"sort"
)

// NodeNameMeta 存储节点名称解析结果。
type NodeNameMeta struct {
	Region    string  // 映射后的地区名
	RegionRaw string  // 原始地区名
	Mult      float64 // 流量倍率（如 1.5, 2.0）
	IsInfo    bool    // 是否为信息节点（非代理节点），为 true 时其他字段无效
}

// NodeBase 从 Clash YAML 中解析出的节点的、原始基础信息。
type NodeBase struct {
	Name     string `json:"name"`     // 节点名称（如 "中国香港-IPLC-C6D-HK-1-流量倍率:1"）
	Type     string `json:"type"`     // 协议类型（如 "trojan"、"vmess"、"vless"、"ss"）
	Server   string `json:"server"`   // 服务器地址
	Port     int    `json:"port"`     // 端口号
	Password string `json:"password"` // 密码
}

func (n *NodeBase) Validate() error {
	if n.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if n.Type == "" {
		return fmt.Errorf("type is empty")
	}
	if n.Server == "" {
		return fmt.Errorf("server is empty")
	}
	if n.Port <= 0 {
		return fmt.Errorf("port is invalid")
	}
	return nil
}

// Node 表示从 Clash YAML 中解析出的一个代理节点。
type Node struct {
	NodeBase
	NodeNameMeta
	raw map[string]any // 原始 Clash YAML 字段（透传给渲染器，避免为每种协议定义结构体）
}

// UID 根据节点信息生成的唯一标识符，排重依据
func (n *Node) UID() string {
	return fmt.Sprintf("%s|%s|%d|%s", n.Type, n.Server, n.Port, n.Password)
}

func (n *Node) SetRaw(r map[string]any) {
	n.raw = r
}

func (n *Node) GetRaw() map[string]any {
	// 替换 raw 中的 name
	if n.Name != "" {
		if _, ok := n.raw["Name"]; ok {
			n.raw["Name"] = n.Name
		}
		if _, ok := n.raw["name"]; ok {
			n.raw["name"] = n.Name
		}
	}
	return n.raw
}

type Nodes struct {
	NormalNodes []*Node // 过滤后的节点列表
	InfoNodes   []*Node // 信息节点
}

func (r *Nodes) NormalNodeNames() []string {
	var result []string
	for _, n := range r.NormalNodes {
		result = append(result, n.Name)
	}
	sort.Strings(result)
	return result
}

func (r *Nodes) NormalNodeRaws() []map[string]any {
	var result []map[string]any
	for _, n := range r.NormalNodes {
		result = append(result, n.GetRaw())
	}
	return result
}

func (r *Nodes) InfoNames() []string {
	var result []string
	for _, n := range r.InfoNodes {
		result = append(result, n.Name)
	}
	sort.Strings(result)
	return result
}

func (r *Nodes) All() []*Node {
	result := make([]*Node, 0, len(r.NormalNodes)+len(r.InfoNodes))
	result = append(result, r.NormalNodes...)
	result = append(result, r.InfoNodes...)
	return result
}
