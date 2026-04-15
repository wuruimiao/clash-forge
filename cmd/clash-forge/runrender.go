package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
	"github.com/wuruimiao/clash-forge/pkg/renderer"
	clashrenderer "github.com/wuruimiao/clash-forge/pkg/renderer/clash"
	mihomorenderer "github.com/wuruimiao/clash-forge/pkg/renderer/mihomo"
	singboxrenderer "github.com/wuruimiao/clash-forge/pkg/renderer/singbox"
)

type RenderCmd struct {
	SharedCommandOptions `embed:""`
	Format               string   `name:"format" short:"f" help:"Output format (${enum})." enum:"singbox,clash,mihomo" env:"FORMAT" required:""`
	Output               string   `name:"output" short:"o" help:"Output file path." env:"OUTPUT"`
	Secret               string   `name:"secret" help:"Dashboard/API secret." env:"SECRET" required:""`
	Auth                 []string `name:"auth" help:"Proxy auth user:pass[@port]." env:"AUTH" sep:","`
	TUN                  bool     `name:"tun" help:"Enable TUN mode." env:"TUN"`
	UTLS                 string   `name:"utls" help:"uTLS fingerprint." env:"UTLS"`
	TFO                  bool     `name:"tfo" help:"Enable TCP Fast Open." env:"TFO"`
	Mux                  bool     `name:"mux" help:"Enable multiplex." env:"MUX"`
	AllowLAN             bool     `name:"allow-lan" default:"true" negatable:"" help:"Allow LAN connections." env:"ALLOW_LAN"`
	ExternalUI           string   `name:"external-ui" help:"External UI path." env:"EXTERNAL_UI"`
}

func parseAuths(entries []string) (map[uint16][]model.Auth, error) {
	result := make(map[uint16][]model.Auth)
	for _, entry := range entries {
		auth, err := parseAuthEntry(entry)
		if err != nil {
			return nil, err
		}
		result[auth.Port] = append(result[auth.Port], auth)
	}
	return result, nil
}

// parseAuthEntry 解析单条认证条目，格式为 user:pass[@port]。
// 密码只允许大小写英文和数字，因此 @ 可安全作为 port 分隔符。
func parseAuthEntry(entry string) (model.Auth, error) {
	// 找最后一个 @ 分隔 port
	var userPass, portStr string
	if lastAt := strings.LastIndex(entry, "@"); lastAt != -1 {
		userPass = entry[:lastAt]
		portStr = entry[lastAt+1:]
	} else {
		userPass = entry
	}

	// 分割 user:pass
	parts := strings.SplitN(userPass, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return model.Auth{}, fmt.Errorf("invalid auth entry (expected user:pass[@port]): %s", entry)
	}

	var port uint16
	if portStr != "" {
		p, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return model.Auth{}, fmt.Errorf("invalid port in auth entry: %s", entry)
		}
		port = uint16(p)
	}

	return model.Auth{
		Username: parts[0],
		Password: parts[1],
		Port:     port,
	}, nil
}

func buildRenderer(format string) (renderer.Renderer, error) {
	switch format {
	case model.SingBox:
		return &singboxrenderer.Renderer{}, nil
	case model.Clash:
		return &clashrenderer.Renderer{}, nil
	case model.MiHoMo:
		return &mihomorenderer.Renderer{}, nil
	default:
		return nil, fmt.Errorf("--format is required (singbox, clash, mihomo)")
	}
}

func (r *RenderCmd) Run() error {
	compiled, err := r.SharedCommandOptions.compile()
	if err != nil {
		return err
	}
	result, err := compiled.executePipeline()
	if err != nil {
		return err
	}

	authUsers, err := parseAuths(r.Auth)
	if err != nil {
		return err
	}

	renderer, err := buildRenderer(r.Format)
	if err != nil {
		return err
	}

	rendered, err := renderer.Render(&model.RenderInput{
		Nodes:      result.filteredNodes,
		NodeGroups: result.groups,
		DNSRules:   result.dnsRules,
		Options: &model.RenderOptions{
			Secret:          r.Secret,
			Auths:           authUsers,
			EnableTUN:       r.TUN,
			UTLSFingerprint: r.UTLS,
			EnableTFO:       r.TFO,
			EnableMux:       r.Mux,
			AllowLAN:        r.AllowLAN,
			ExternalUI:      r.ExternalUI,
			Debug:           r.Debug,
			SkipAuths:       r.SkipAuths,
		},
	})
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if r.Output == "" {
		_, err = os.Stdout.Write(rendered)
		return err
	}
	return util.WriteFileAtomic(r.Output, rendered)
}
