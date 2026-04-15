package model

// RenderOptions 存储通过 CLI 标志传入的渲染选项。
type RenderOptions struct {
	Secret          string     // API/dashboard 密钥
	AuthUsers       []AuthUser // 代理端口认证用户列表
	EnableTUN       bool       // 是否启用 TUN 模式（透明代理）
	UTLSFingerprint string     // uTLS 指纹伪装（如 chrome、firefox、safari）
	EnableTFO       bool       // 是否启用 TCP Fast Open
	EnableMux       bool       // 是否启用多路复用
	AllowLAN        bool       // 是否允许局域网连接
	ExternalUI      string     // 外部 UI 路径（Clash 专用）
	Debug           bool       // 是否启用调试模式
}

// AuthUser 表示代理认证的用户名和密码对。
type AuthUser struct {
	Username string
	Password string
}

// RenderInput 是传递给渲染器的完整输入，包含节点、分组、原始配置和渲染选项。
type RenderInput struct {
	*Nodes
	*NodeGroups
	DNSRules *DNSRules      // 源配置中保留的 DNS 和规则
	Options  *RenderOptions // 用户指定的渲染选项
}
