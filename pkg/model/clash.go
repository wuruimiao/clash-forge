package model

import "fmt"

// DNSRules 保存源 Clash 配置中的非代理部分，用于渲染时还原。
type DNSRules struct {
	DNS   map[string]any `yaml:"dns"`   // 原始 DNS 配置
	Rules []string       `yaml:"rules"` // 原始路由规则
}

func (d *DNSRules) Validate() error {
	if len(d.DNS) == 0 {
		return fmt.Errorf("dns is empty")
	}
	if len(d.Rules) == 0 {
		return fmt.Errorf("rules is empty")
	}
	return nil
}

// ClashConfig 表示 Clash YAML 配置的顶层结构，只提取 proxies、dns、rules 三个关键字段。
type ClashConfig struct {
	Proxies  []map[string]any `yaml:"proxies"`
	DNSRules `yaml:",inline"`
}
