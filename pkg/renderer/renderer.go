// Package renderer 定义了配置渲染器的通用接口。
// sing-box 和 Clash/Mihomo 分别实现此接口，将统一的节点数据转为各自格式的配置文件。
package renderer

import "github.com/wuruimiao/clash-forge/pkg/model"

// Renderer 是配置文件生成器的接口。
// 每个实现负责将 RenderInput（节点、分组、选项）转化为目标代理软件的完整配置。
type Renderer interface {
	Render(input *model.RenderInput) ([]byte, error) // 生成配置文件内容
	Name() string                                     // 返回渲染器名称（如 "singbox"、"clash"）
}
