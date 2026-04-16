package singbox

import (
	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

// buildTLS 构建 TLS 配置。
// trojan 和 hysteria2 默认启用 TLS，其他协议需要在 Raw 中有 tls: true。
// 支持 SNI、跳过证书验证、ALPN、uTLS 指纹伪装和 Reality 配置。
func buildTLS(node *model.Node, opts *model.RenderOptions) map[string]any {
	hasTLS := false
	switch node.Type {
	case "trojan", "hysteria2":
		hasTLS = true // 这两种协议强制启用 TLS
	default:
		if v, ok := node.GetRaw()["tls"]; ok {
			hasTLS, _ = v.(bool)
		}
	}

	if !hasTLS {
		return nil
	}

	tls := map[string]any{"enabled": true}

	// 服务器名称指示（SNI），Clash 配置中可能使用 "sni" 或 "servername"
	if sni := util.StrVal(node.GetRaw(), "sni"); sni != "" {
		tls["server_name"] = sni
	} else if sn := util.StrVal(node.GetRaw(), "servername"); sn != "" {
		tls["server_name"] = sn
	}

	// 跳过证书验证（不安全，仅用于自签名证书场景）
	if v, ok := node.GetRaw()["skip-cert-verify"]; ok {
		if b, ok := v.(bool); ok && b {
			tls["insecure"] = true
		}
	}

	// ALPN 协议列表
	if alpn, ok := node.GetRaw()["alpn"]; ok {
		tls["alpn"] = alpn
	}

	// uTLS 指纹伪装：模拟特定浏览器的 TLS 指纹，避免被识别
	if opts.UTLSFingerprint != "" {
		tls["utls"] = map[string]any{
			"enabled":     true,
			"fingerprint": opts.UTLSFingerprint,
		}
	}

	// Reality 配置：VLESS Reality 协议使用的特殊 TLS 握手
	if realityOpts, ok := node.GetRaw()["reality-opts"].(map[string]any); ok {
		tls["reality"] = buildRealityFields(realityOpts)
	}

	return tls
}

// buildRealityFields 构建 Reality 配置子对象。
func buildRealityFields(raw map[string]any) map[string]any {
	reality := map[string]any{"enabled": true}
	if pk := util.StrVal(raw, "public-key"); pk != "" {
		reality["public_key"] = pk
	}
	if sid := util.StrVal(raw, "short-id"); sid != "" {
		reality["short_id"] = sid
	}
	return reality
}

// buildTransport 构建传输层配置。
// 根据 Clash 的 network 字段映射为 sing-box 的传输层类型：ws → WebSocket, grpc → gRPC, http/h2 → HTTP。
// TCP 直连或空值时不需要额外传输层配置。
func buildTransport(node *model.Node) map[string]any {
	network := util.StrVal(node.GetRaw(), "network")
	if network == "" || network == "tcp" {
		return nil
	}

	switch network {
	case "ws":
		// WebSocket 传输
		transport := map[string]any{"type": "ws"}
		if wsOpts, ok := node.GetRaw()["ws-opts"].(map[string]any); ok {
			if path := util.StrVal(wsOpts, "path"); path != "" {
				transport["path"] = path
			}
			if headers, ok := wsOpts["headers"].(map[string]any); ok {
				transport["headers"] = headers
			}
		}
		return transport
	case "grpc":
		// gRPC 传输
		transport := map[string]any{"type": "grpc"}
		if grpcOpts, ok := node.GetRaw()["grpc-opts"].(map[string]any); ok {
			if sn := util.StrVal(grpcOpts, "grpc-service-name"); sn != "" {
				transport["service_name"] = sn
			}
		}
		return transport
	case "http", "h2":
		// HTTP/2 传输
		transport := map[string]any{"type": "http"}
		if httpOpts, ok := node.GetRaw()["http-opts"].(map[string]any); ok {
			if path, ok := httpOpts["path"]; ok {
				transport["path"] = path
			}
			if headers, ok := httpOpts["headers"].(map[string]any); ok {
				transport["headers"] = headers
			}
		}
		return transport
	}
	return nil
}
