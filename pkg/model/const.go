package model

const (
	All      = "auto"
	AutoMult = "auto-mult"
	Proxy    = "Proxy"

	SingBox = "singbox"
	Clash   = "clash"
	MiHoMo  = "mihomo"

	// `^(?P<region>.*?)-(?P<id>.*?)-流量倍率:(?P<mult>\d+(\.\d+)?)$`,
	DefaultNamePattern = `^(?P<region>.*?)-.*-流量倍率:(?P<mult>\d+(\.\d+)?)$`
	DefaultTestUrl     = "https://www.gstatic.com/generate_204"

	UrlTestIntervalSec = 300
	UrlTestIntervalMin = "3m"

	ToolConfigName = "clash-forge.yaml"
)
