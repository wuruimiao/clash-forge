package singbox

import "github.com/wuruimiao/clash-forge/pkg/model"

// buildRoute 构建路由规则段。
// 路由策略：嗅探协议 → 劫持 DNS → 私有 IP 直连 → 中国站点/IP 直连 → 其余走 Proxy 代理组。
// 使用远程规则集（geosite-cn、geoip-cn）实现中国流量直连分流。
func buildRoute() map[string]any {
	return map[string]any{
		"rules": []any{
			map[string]any{"action": "sniff"},                                                 // 协议嗅探
			map[string]any{"protocol": "dns", "action": "hijack-dns"},                         // DNS 劫持
			map[string]any{"ip_is_private": true, "action": "route", "outbound": "direct"},    // 私有 IP 直连
			map[string]any{"rule_set": "geosite-cn", "action": "route", "outbound": "direct"}, // 中国站点直连
			map[string]any{"rule_set": "geoip-cn", "action": "route", "outbound": "direct"},   // 中国 IP 直连
		},
		"rule_set": []any{
			map[string]any{
				"type": "remote", "tag": "geosite-cn", "format": "binary",
				"url":             "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs",
				"download_detour": "direct",
			},
			map[string]any{
				"type": "remote", "tag": "geoip-cn", "format": "binary",
				"url":             "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs",
				"download_detour": "direct",
			},
		},
		"final":                 "Proxy", // 未匹配规则的流量走 Proxy 代理组
		"auto_detect_interface": true,    // 自动检测出站网络接口
	}
}

// buildExperimental 构建实验性功能配置段。
// 包含缓存文件（持久化节点测速结果）和 Clash API（供外部 UI 管理面板使用）。
func buildExperimental(opts *model.RenderOptions) map[string]any {
	clashAPI := map[string]any{
		"external_controller": "0.0.0.0:9090",
	}
	if opts.Secret != "" {
		clashAPI["secret"] = opts.Secret
	}

	return map[string]any{
		"cache_file": map[string]any{
			"enabled":      true,
			"store_fakeip": true,
		},
		"clash_api": clashAPI,
	}
}
