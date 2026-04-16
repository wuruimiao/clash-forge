package mihomo

import (
	"fmt"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
	clashrenderer "github.com/wuruimiao/clash-forge/pkg/renderer/clash"
	"gopkg.in/yaml.v3"
)

type Renderer struct {
}

func (r *Renderer) Name() string {
	return model.MiHoMo
}
func (r *Renderer) Render(input *model.RenderInput) ([]byte, error) {
	clash := &clashrenderer.Renderer{}
	cfg := util.NewYamlOrderedMap()

	cfg.Add("mode", "rule")

	clash.BuildBaseConfig(cfg, input)

	buildConfig(cfg, input)

	clash.BuildNet(cfg, input)
	return yaml.Marshal(cfg.Node)
}

func buildConfig(cfg *util.YamlOrderedMap, input *model.RenderInput) {
	// ui
	if input.Options.ExternalUI != "" {
		cfg.Add("external-ui", input.Options.ExternalUI)
	} else {
		cfg.Add("external-ui", "/root/.config/mihomo/public")
	}
	// 默认：https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip
	cfg.Add("external-ui-url", "https://github.com/haishanh/yacd/archive/gh-pages.zip")

	cfg.Add("geo-auto-update", true)
	cfg.Add("geo-update-interval", 24)
	cfg.Add("geodata-mode", false)
	cfg.Add("geox-url", map[string]string{
		"geoip":   "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geoip.dat",
		"geosite": "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geosite.dat",
		"mmdb":    "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/country.mmdb",
		"asn":     "https://github.com/xishang0128/geoip/releases/download/latest/GeoLite2-ASN.mmdb",
	})

	// 统一延迟
	cfg.Add("unified-delay", true)

	// 鉴权
	var listeners []*yaml.Node
	for port, users := range input.Options.Auths {
		if port == 0 {
			// 已由 clash 处理
			continue
		}
		listener := util.NewYamlOrderedMap()
		listener.Add("name", fmt.Sprintf("mixed-%d", port))
		listener.Add("type", "mixed")
		listener.Add("port", port)
		listener.Add("listen", "0.0.0.0")
		var listenerUsers []map[string]string
		for _, u := range users {
			listenerUsers = append(listenerUsers, map[string]string{
				"username": u.Username,
				"password": u.Password,
			})
		}
		listener.Add("users", listenerUsers)
		listeners = append(listeners, listener.Node)
	}

	if len(listeners) > 0 {
		cfg.Add("listeners", listeners)
	}

	// 跳过鉴权的ip段
	if len(input.Options.SkipAuths) > 0 {
		cfg.Add("skip-auth-prefixes", input.Options.SkipAuths)
	}

}
