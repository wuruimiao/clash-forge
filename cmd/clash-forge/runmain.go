package main

import (
	"fmt"

	"github.com/wuruimiao/clash-forge/pkg/grouper"
	"github.com/wuruimiao/clash-forge/pkg/model"
	"github.com/wuruimiao/clash-forge/pkg/parser"
)

type pipelineResult struct {
	allNodes      []*model.Node
	filteredNodes *model.Nodes
	groups        *model.NodeGroups
	dnsRules      *model.DNSRules
}

// executePipeline 主流程
func (r *compiledSharedConfig) executePipeline() (*pipelineResult, error) {
	// 节点名解析器
	nodeNameParser, err := (&parser.NodeNameParserOptions{
		Re:           r.NameRe,
		ExtraRegions: r.RegionMap,
	}).New()
	if err != nil {
		return nil, err
	}

	// 获取 input 内容并解析
	parsedConfig, err := handleSources(r.Inputs, nodeNameParser)
	if err != nil {
		return nil, err
	}

	// 获取所有节点的copy列表
	filteredNodes := parsedConfig.GetNodes()
	// 过滤节点
	filteredNodes.NormalNodes = r.FilterSpec.Apply(filteredNodes.NormalNodes)
	if len(filteredNodes.NormalNodes) == 0 {
		return nil, fmt.Errorf("normal nodes filteredNodes out; try relaxing filter criteria")
	}

	// 节点分组
	groups := grouper.Group(filteredNodes.NormalNodes)

	return &pipelineResult{
		allNodes:      parsedConfig.Nodes,
		filteredNodes: filteredNodes,
		groups:        groups,
		dnsRules:      parsedConfig.DNSRules,
	}, nil
}
