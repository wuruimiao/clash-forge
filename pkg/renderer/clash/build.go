package clash

import "github.com/wuruimiao/clash-forge/pkg/model"

// buildUrlTest 构建一个 URL 测试代理组。自动测速，每 x 秒测试一次
func buildUrlTest(name string, proxies []string) map[string]any {
	return map[string]any{
		"name":     name,
		"type":     "url-test",
		"url":      model.DefaultTestUrl,
		"interval": model.UrlTestIntervalSec,
		"proxies":  proxies,
	}
}

// buildSelect 构建一个选择代理组。
func buildSelect(name string, proxies []string) map[string]any {
	return map[string]any{
		"name":    name,
		"type":    "select",
		"proxies": proxies,
	}
}

// buildFallback 构建一个故障转移代理组。自动测速，每 x 秒测试一次
func buildFallback(name string, proxies []string) map[string]any {
	return map[string]any{
		"name":     name,
		"type":     "fallback",
		"url":      model.DefaultTestUrl,
		"interval": model.UrlTestIntervalSec,
		"proxies":  proxies,
	}
}

// buildDNS 返回默认 DNS 配置。
// 主 DNS：阿里 223.5.5.5 和腾讯 119.29.29.29（国内快速解析）。
// 备用 DNS：Google 8.8.8.8 和 Cloudflare 1.0.0.1（国外域名兜底）。
func buildDNS() map[string]any {
	return map[string]any{
		"enable": true,
		"nameserver": []string{
			"223.5.5.5",
			"119.29.29.29",
		},
		"fallback": []string{
			"8.8.8.8",
			"tls://1.0.0.1:853",
		},
	}
}

// buildRules 返回默认路由规则。
// 中国 IP 直连，其余流量走 Proxy 代理组。
func buildRules() []string {
	return []string{
		"GEOIP,CN,DIRECT",
		"MATCH,Proxy",
	}
}
