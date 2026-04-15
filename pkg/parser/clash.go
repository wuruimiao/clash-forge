// Package parser 负责解析 Clash 格式的 YAML 配置文件，
// 将其中的代理节点提取为统一的 model.Node 结构。
package parser

import (
	"fmt"
	"log/slog"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
	"gopkg.in/yaml.v3"
)

// ParseClashYAML 解析 Clash YAML 数据，返回节点列表和原始配置。
// nodeNameParser 为 nil 时使用默认解析逻辑。
func ParseClashYAML(data []byte, nodeNameParser *NodeNameParser) (*model.ParsedConfig, error) {
	var cfg model.ClashConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse clash yaml: %w", err)
	}

	var nodes []*model.Node
	for _, raw := range cfg.Proxies {
		node, ok := parseProxy(raw, nodeNameParser)
		if !ok {
			continue
		}
		nodes = append(nodes, node)
	}

	parsed := &model.ParsedConfig{
		Nodes: nodes,
	}
	dnsRules := &model.DNSRules{
		DNS:   cfg.DNS,
		Rules: cfg.Rules,
	}
	if err := dnsRules.Validate(); err != nil {
		slog.Warn("invalid dns rules, will ignore.", "err", err)
	} else {
		// 保留源配置中的 DNS 和路由规则，供渲染器还原使用
		parsed.DNSRules = dnsRules
	}
	return parsed, nil
}

// parseProxy 解析单个代理条目。
// 首先通过节点名称判断是否为信息节点（到期/流量通知等），
// 然后校验必填字段（name、type、server、port），缺失则跳过并记录警告。
func parseProxy(raw map[string]any, nodeNameParser *NodeNameParser) (*model.Node, bool) {
	nodeBase, err := util.Map2Struct[model.NodeBase](raw)
	if err != nil {
		slog.Error("skipping proxy: invalid node base", "err=", err, "raw=", raw)
		return nil, false
	}
	if err := nodeBase.Validate(); err != nil {
		slog.Error("skipping proxy: validation failed", "err=", err, "name=", nodeBase.Name)
		return nil, false
	}

	// 通过节点名称解析地区、标识符和倍率
	nodeNameMeta := nodeNameParser.Parse(nodeBase.Name)

	node := &model.Node{
		NodeBase:     nodeBase,
		NodeNameMeta: nodeNameMeta,
	}
	node.SetRaw(raw)
	return node, true
}
