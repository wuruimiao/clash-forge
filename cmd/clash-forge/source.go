package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/wuruimiao/clash-forge/pkg/model"
	"github.com/wuruimiao/clash-forge/pkg/parser"
	"github.com/wuruimiao/clash-forge/pkg/source"
)

func handleSources(inputs []string, nodeNameParser *parser.NodeNameParser) (*model.ParsedConfig, error) {
	dataMap, fetchErrs := fetchSources(inputs)
	result, parseErrors := parseClashSources(dataMap, nodeNameParser)
	if len(result.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes from any source; errors: %s", strings.Join(append(fetchErrs, parseErrors...), "; "))
	}
	return result, nil
}

// fetchSources 遍历所有输入源，获取原始数据。
// 返回输入到数据的映射，以及获取失败的错误列表。
// 单个源失败不会中断流程，仅记录警告。
func fetchSources(inputs []string) (map[string][]byte, []string) {
	dataMap := make(map[string][]byte)
	var errors []string

	for _, input := range inputs {
		data, err := source.FetchOne(input)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", input, err))
			slog.Warn("failed to fetch source", "input", input, "error", err)
			continue
		}
		dataMap[input] = data
	}
	return dataMap, errors
}

// parseClashSources 解析已获取的原始数据。
// 多源合并：合并节点列表，只保留第一个成功源的 DNS 和规则配置。
// 单个源解析失败不会中断流程，仅记录警告。
func parseClashSources(dataMap map[string][]byte, nodeNameParser *parser.NodeNameParser) (*model.ParsedConfig, []string) {
	var errors []string
	parsedConfig := &model.ParsedConfig{}

	for input, data := range dataMap {
		parsed, err := parser.ParseClashYAML(data, nodeNameParser)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", input, err))
			slog.Warn("failed to parse source", "input", input, "error", err)
			continue
		}

		parsedConfig.Nodes = append(parsedConfig.Nodes, parsed.Nodes...)

		if parsedConfig.DNSRules == nil {
			parsedConfig.DNSRules = parsed.DNSRules
		}
	}

	return parsedConfig, errors
}
