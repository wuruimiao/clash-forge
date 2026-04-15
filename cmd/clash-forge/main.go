package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/wuruimiao/clash-forge/internal/util"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	// 构建 kong 文件 opts
	args := os.Args[1:]
	configPath := resolveConfigPath(args)

	cli := newCLI()
	parser, err := cli.buildParser(configPath)
	if err != nil {
		return err
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		parser.FatalIfErrorf(err)
		return err
	}

	if v := ctx.Selected().Target.FieldByName("Debug"); v.IsValid() && v.Bool() {
		slog.Debug(strings.Repeat("=", 105))
		slog.Debug("Config snapshot:")
		util.JsonV(ctx.Selected().Target.Interface())
		slog.Debug(strings.Repeat("=", 105))
	}

	// 运行
	if err := ctx.Run(); err != nil {
		return err
	}

	return nil
}
