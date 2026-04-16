package singbox

import "github.com/wuruimiao/clash-forge/pkg/model"

// buildDNS 构建 DNS 配置段。
// 策略：国外域名走远程 DNS（Google 8.8.8.8 over TLS），中国域名走本地 DNS（阿里 223.5.5.5 UDP）。
func buildDNS() map[string]any {
	return map[string]any{
		"servers": []any{
			map[string]any{"type": "tls", "tag": "dns-remote", "server": "8.8.8.8"},
			map[string]any{"type": "udp", "tag": "dns-local", "server": "223.5.5.5"},
		},
		"rules": []any{
			map[string]any{"outbound": "any", "server": "dns-local"},                // 直连出站使用本地 DNS
			map[string]any{"domain_suffix": []string{".cn"}, "server": "dns-local"}, // .cn 域名使用本地 DNS
		},
		"final": "dns-remote", // 其余域名使用远程 DNS
	}
}

func buildSelector(name string, tags []string) map[string]any {
	return map[string]any{
		"type":      "selector",
		"tag":       name,
		"outbounds": tags,
		"default":   "auto",
	}
}

func buildUrlTest(name string, tags []string) map[string]any {
	return map[string]any{
		"type":      "urltest",
		"tag":       name,
		"outbounds": tags,
		"url":       model.DefaultTestUrl,
		"interval":  model.UrlTestIntervalMin,
	}
}
