package singbox

import (
	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

// convertOutbound 将 Clash 格式的 Node 转换为 sing-box 出站配置。
// 支持的协议：trojan、vmess、vless、ss（shadowsocks）、hysteria2、tuic。
// 不支持的协议（如 ssr）返回 nil。
func convertOutbound(node *model.Node, opts *model.RenderOptions) map[string]any {
	// 所有协议共有的基础字段
	out := map[string]any{
		"tag":         node.Name,
		"server":      node.Server,
		"server_port": node.Port,
	}

	// 根据协议类型设置特定字段
	var protoFields map[string]any
	switch node.Type {
	case "trojan":
		out["type"] = "trojan"
		protoFields = convertTrojan(node.GetRaw())
	case "vmess":
		out["type"] = "vmess"
		protoFields = convertVMess(node.GetRaw())
	case "vless":
		out["type"] = "vless"
		protoFields = convertVLESS(node.GetRaw())
	case "ss":
		// Clash 的 "ss" 映射为 sing-box 的 "shadowsocks"
		out["type"] = "shadowsocks"
		protoFields = convertShadowsocks(node.GetRaw())
	case "hysteria2":
		out["type"] = "hysteria2"
		protoFields = convertHysteria2(node.GetRaw())
	case "tuic":
		out["type"] = "tuic"
		protoFields = convertTUIC(node.GetRaw())
	case "ssr":
		// SSR 协议不受 sing-box 支持，返回 nil
		return nil
	default:
		// 未知协议直接透传类型名
		out["type"] = node.Type
	}

	// 合并协议特定字段
	for k, v := range protoFields {
		out[k] = v
	}

	// 构建 TLS 配置（包括 Reality）
	tls := buildTLS(node, opts)
	if tls != nil {
		out["tls"] = tls
	}

	// 构建传输层配置（WebSocket、gRPC、HTTP）
	transport := buildTransport(node)
	if transport != nil {
		out["transport"] = transport
	}

	// 多路复用：多个连接共享同一底层通道，减少握手开销
	if opts.EnableMux {
		out["multiplex"] = map[string]any{"enabled": true}
	}

	// TCP Fast Open：减少 TCP 连接建立的 RTT
	if opts.EnableTFO {
		out["tcp_fast_open"] = true
	}

	return out
}

// convertTrojan 构建 trojan 协议特定字段。
func convertTrojan(raw map[string]any) map[string]any {
	return map[string]any{
		"password": util.StrVal(raw, "password"),
	}
}

// convertVMess 构建 vmess 协议特定字段。
func convertVMess(raw map[string]any) map[string]any {
	out := map[string]any{
		"uuid":     util.StrVal(raw, "uuid"),
		"security": util.StrValOr(raw, "cipher", "auto"),
	}
	if alterId, ok := raw["alterId"]; ok {
		out["alter_id"] = alterId
	}
	return out
}

// convertVLESS 构建 vless 协议特定字段。
func convertVLESS(raw map[string]any) map[string]any {
	out := map[string]any{
		"uuid": util.StrVal(raw, "uuid"),
	}
	if flow := util.StrVal(raw, "flow"); flow != "" {
		out["flow"] = flow // 流控模式（如 xtls-rprx-vision）
	}
	return out
}

// convertShadowsocks 构建 shadowsocks 协议特定字段。
func convertShadowsocks(raw map[string]any) map[string]any {
	return map[string]any{
		"method":   util.StrVal(raw, "cipher"),
		"password": util.StrVal(raw, "password"),
	}
}

// convertHysteria2 构建 hysteria2 协议特定字段。
func convertHysteria2(raw map[string]any) map[string]any {
	out := map[string]any{
		"password": util.StrVal(raw, "password"),
	}
	if up := util.StrVal(raw, "up"); up != "" {
		out["up_mbps"] = raw["up"]
	}
	if down := util.StrVal(raw, "down"); down != "" {
		out["down_mbps"] = raw["down"]
	}
	// 混淆配置（可选）
	if obfs := util.StrVal(raw, "obfs"); obfs != "" {
		out["obfs"] = map[string]any{
			"type":     obfs,
			"password": util.StrVal(raw, "obfs-password"),
		}
	}
	return out
}

// convertTUIC 构建 tuic 协议特定字段。
func convertTUIC(raw map[string]any) map[string]any {
	out := map[string]any{
		"uuid":     util.StrVal(raw, "uuid"),
		"password": util.StrVal(raw, "password"),
	}
	if cc := util.StrVal(raw, "congestion-controller"); cc != "" {
		out["congestion_control"] = cc // 拥塞控制算法（如 bbr、cubic）
	}
	if urm := util.StrVal(raw, "udp-relay-mode"); urm != "" {
		out["udp_relay_mode"] = urm // UDP 中继模式（如 native、quic）
	}
	return out
}
