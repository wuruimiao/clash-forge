// Package singbox 实现了 sing-box 格式的 JSON 配置渲染器。
// 生成包含 log、dns、inbounds、outbounds、route、experimental 的完整 sing-box 配置。
package singbox

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/wuruimiao/clash-forge/pkg/model"
)

// Renderer 是 sing-box JSON 配置的渲染器。
type Renderer struct{}

// Name 返回渲染器名称。
func (r *Renderer) Name() string { return "singbox" }

// Render 生成完整的 sing-box JSON 配置。
// 配置包含六个顶层段落：日志、DNS、入站、出站、路由和实验性功能。
func (r *Renderer) Render(input *model.RenderInput) ([]byte, error) {
	cfg := map[string]any{
		"log":          buildLog(),
		"dns":          buildDNS(),
		"inbounds":     buildInbounds(input.Options),
		"outbounds":    buildOutbounds(input.Nodes, input.NodeGroups, input.Options),
		"route":        buildRoute(),
		"experimental": buildExperimental(input.Options),
	}

	return json.MarshalIndent(cfg, "", "  ")
}

// buildLog 构建日志配置段。
func buildLog() map[string]any {
	return map[string]any{"level": "info", "timestamp": true}
}

// buildInbounds 构建入站配置段。
// 可选 TUN 入站（透明代理）+ 必选 mixed 入站（HTTP/SOCKS5 混合代理）。
// Port == 0 的用户放入默认 mixed-in，Port != 0 的用户生成独立的 mixed inbound。
func buildInbounds(opts *model.RenderOptions) []any {
	var inbounds []any

	// TUN 入站：透明代理模式，需要 root 权限
	if opts.EnableTUN {
		inbounds = append(inbounds, map[string]any{
			"type":         "tun",
			"tag":          "tun-in",
			"address":      []string{"172.18.0.1/30", "fdfe:dcba:9876::1/126"},
			"mtu":          9000,
			"auto_route":   true,
			"strict_route": true,
			"stack":        "mixed",
		})
	}

	// 默认 mixed 入站：HTTP + SOCKS5 混合代理端口
	mixed := map[string]any{
		"type":        "mixed",
		"tag":         "mixed-in",
		"listen":      "0.0.0.0",
		"listen_port": 2080,
	}

	for port, users := range opts.Auths {
		if port == 0 {
			var us []any
			for _, u := range users {
				us = append(us, map[string]any{
					"username": u.Username,
					"password": u.Password,
				})
			}
			mixed["users"] = us
		} else {
			inbound := map[string]any{
				"type":        "mixed",
				"tag":         fmt.Sprintf("mixed-in-%d", port),
				"listen":      "0.0.0.0",
				"listen_port": port,
			}
			var inboundUsers []any
			for _, u := range users {
				inboundUsers = append(inboundUsers, map[string]any{
					"username": u.Username,
					"password": u.Password,
				})
			}
			inbound["users"] = inboundUsers
			inbounds = append(inbounds, inbound)
		}
	}
	inbounds = append(inbounds, mixed)

	return inbounds
}

// buildOutbounds 构建出站配置段。
// 结构：分组出站（selector/urltest）→ auto 全局测速组 → 各节点出站 → 特殊出站（direct/block/dns）。
func buildOutbounds(
	allNodes *model.Nodes,
	groups *model.NodeGroups,
	opts *model.RenderOptions) []any {
	var outbounds []any

	// 收集所有节点名称作为 tag
	allTags := allNodes.NormalNodeNames()

	// TODO: 增加 model.AutoMult
	proxyTags := []string{model.All}
	groups.GroupNodeNames().Range(func(groupName string, nodeNames []string) {
		// 地区/倍率组（urltest 类型）：自动测速选最快节点
		outbounds = append(outbounds, buildUrlTest(groupName, nodeNames))
		// Proxy 全局组（selector 类型）：包含 auto 和所有子分组，供用户手动切换
		proxyTags = append(proxyTags, groupName)
	})
	outbounds = append([]any{buildSelector("Proxy", proxyTags)}, outbounds...)

	// auto 全局自动测速组：包含所有节点
	outbounds = append(outbounds, buildUrlTest(model.All, allTags))

	// 各节点的独立出站（协议转换由 convertOutbound 完成）
	for _, n := range allNodes.NormalNodes {
		out := convertOutbound(n, opts)
		if out == nil {
			slog.Warn("skipping unsupported protocol for sing-box", "name", n.Name, "type", n.Type)
			continue
		}
		outbounds = append(outbounds, out)
	}

	// 特殊出站：直连、拦截、DNS
	outbounds = append(outbounds, buildSpecialOutbounds()...)

	return outbounds
}

// buildSpecialOutbounds 返回特殊出站配置（直连、拦截）。
func buildSpecialOutbounds() []any {
	return []any{
		map[string]any{"type": "direct", "tag": "direct"},
		map[string]any{"type": "block", "tag": "block"},
	}
}
