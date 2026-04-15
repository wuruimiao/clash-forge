package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

type CLI struct {
	Render  RenderCmd  `cmd:"" help:"Generate configuration output."`
	Preview PreviewCmd `cmd:"" help:"Preview parsed nodes and grouping."`
	List    ListCmd    `cmd:"" help:"List filtered nodes."`
	Version VersionCmd `cmd:"" help:"Show version."`
}

func resolveConfigPath(args []string) string {
	for i := range args {
		if args[i] == "--config" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return model.ToolConfigName
}

func newCLI() *CLI { return &CLI{} }

func (cli *CLI) buildParser(configPath string) (*kong.Kong, error) {
	opts := []kong.Option{
		kong.Name("clash-forge"),
		kong.UsageOnError(),
	}

	if _, err := os.Stat(configPath); err == nil {
		opts = append(opts, kong.Configuration(util.YamlConfigLoader, configPath))
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	return kong.New(cli, opts...)
}
