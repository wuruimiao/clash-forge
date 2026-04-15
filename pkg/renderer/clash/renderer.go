// Package clash 实现了 Clash/Mihomo 格式的 YAML 配置渲染器。
// 输出兼容 Clash Premium 和 Mihomo（原 Clash.Meta）。
package clash

import (
	"fmt"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
	"gopkg.in/yaml.v3"
)

// Renderer 是 Clash/Mihomo YAML 配置的渲染器。
type Renderer struct {
}

// Name 返回渲染器名称。
func (r *Renderer) Name() string {
	return model.Clash
}

// Render 生成完整的 Clash/Mihomo YAML 配置。
func (r *Renderer) Render(input *model.RenderInput) ([]byte, error) {
	cfg := util.NewYamlOrderedMap()

	cfg.Add("mode", "Rule")

	r.BuildBaseConfig(cfg, input)

	// ui
	if input.Options.ExternalUI != "" {
		cfg.Add("external-ui", input.Options.ExternalUI)
	} else {
		cfg.Add("external-ui", "/root/.config/clash/public")
	}

	r.BuildNet(cfg, input)
	return yaml.Marshal(cfg.Node)
}

// BuildBaseConfig 构建完整的 Clash 配置 map。
// 配置结构：基础设置 → DNS → 代理列表 → 代理分组 → 路由规则。
func (r *Renderer) BuildBaseConfig(cfg *util.YamlOrderedMap, input *model.RenderInput) {
	cfg.Add("mixed-port", 7890)
	cfg.Add("allow-lan", input.Options.AllowLAN)

	if input.Options.Debug {
		cfg.Add("log-level", "debug")
	} else {
		cfg.Add("log-level", "silent")
	}

	// ui密钥
	if input.Options.Secret != "" {
		cfg.Add("secret", input.Options.Secret)
	}

	// ui
	cfg.Add("external-controller", "0.0.0.0:9090")

	// 全局鉴权
	var auth []string
	for _, u := range input.Options.Auths[0] {
		auth = append(auth, fmt.Sprintf("%s:%s", u.Username, u.Password))
	}
	if len(auth) > 0 {
		cfg.Add("authentication", auth)
	}
}

func (r *Renderer) BuildNet(cfg *util.YamlOrderedMap, input *model.RenderInput) {
	// DNS 配置：优先使用源配置中的 DNS，否则使用默认值
	if input.DNSRules != nil && input.DNSRules.DNS != nil {
		cfg.Add("dns", input.DNSRules.DNS)
	} else {
		cfg.Add("dns", buildDNS())
	}

	// 代理列表：直接透传 Raw 字段，保留源配置中的所有协议细节
	cfg.Add("proxies", input.Nodes.NormalNodeRaws())

	// 代理分组
	cfg.Add("proxy-groups", buildProxyGroups(input.Nodes, input.NodeGroups))

	// 路由规则：优先使用源配置中的规则，否则使用默认值
	if input.DNSRules != nil && len(input.DNSRules.Rules) > 0 {
		cfg.Add("rules", input.DNSRules.Rules)
	} else {
		cfg.Add("rules", buildRules())
	}
}

// buildProxyGroups 构建代理分组配置。
// 结构：Proxy（select，包含 Auto + 各子分组名）→ 各子分组（url-test）→ Auto（url-test，所有节点）。
func buildProxyGroups(allNodes *model.Nodes, groups *model.NodeGroups) []map[string]any {
	var result []map[string]any

	// 信息节点
	// for _, n := range allNodes.InfoNames() {
	// 	result = append(result, buildSelect(n, []string{n}))
	// }

	var result2 []map[string]any

	// All 分组 所有可用节点，方便节点测试
	result2 = append(result2, buildUrlTest(model.All, allNodes.NormalNodeNames()))

	// Auto-Mult 分组，分组间 fallback
	result2 = append(result2, buildFallback(model.AutoMult, groups.GroupMultNames()))

	// 所有可用节点和节点分组
	proxyNames := []string{model.All, model.AutoMult}

	groups.GroupNodeNames().Range(func(groupName string, nodeNames []string) {
		// Proxy 组加入这些分组
		proxyNames = append(proxyNames, groupName)
		// 各分组，包括地区/倍率/...分组
		result2 = append(result2, buildUrlTest(groupName, nodeNames))
	})

	proxyNames = append(proxyNames, allNodes.NormalNodeNames()...)

	// 保证 Proxy 在信息节点之前
	result = append(result, buildSelect(model.Proxy, proxyNames))
	result = append(result, result2...)
	return result
}
