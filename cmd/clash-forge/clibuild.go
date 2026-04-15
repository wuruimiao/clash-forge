package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/filter"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

type SharedDirectCommandOptions struct {
	Inputs     []string `name:"input" short:"i" help:"Config source path or URL." env:"INPUT" sep:"," required:""`
	Debug      bool     `name:"debug" help:"Enable debug output and print config snapshot." env:"DEBUG"`
	ConfigFile string   `name:"config" default:"clash-forge.yaml" help:"YAML config file path."` // TODO: 这个可能不需要
	SkipAuths  []string `name:"skip-auth" help:"IP ranges that skip auth verification." env:"SKIP_AUTH" sep:","`
}

type SharedCommandOptions struct {
	SharedDirectCommandOptions `embed:""`
	NamePattern                string   `name:"name-pattern" help:"Regex with named groups region and mult." env:"NAME_PATTERN"`
	RegionMap                  []string `name:"region-map" help:"Extra region mapping entries like 中国:中 or 中国: or 中国. All str started with 中国 will be mapped." env:"REGION_MAP" sep:","`
	IncludeRegion              []string `name:"include-region" help:"Keep only these regions." env:"INCLUDE_REGION" sep:","`
	ExcludeRegion              []string `name:"exclude-region" help:"Exclude these regions." env:"EXCLUDE_REGION" sep:","`
	ExcludeType                []string `name:"exclude-type" help:"Exclude these node types." env:"EXCLUDE_TYPE" sep:","`
	ExcludePattern             []string `name:"exclude-pattern" help:"Exclude node names by regex." env:"EXCLUDE_PATTERN" sep:","`
	MaxMult                    float64  `name:"max-mult" help:"Maximum traffic mult (0 means unlimited)." env:"MAX_MULT"`
	MinMult                    float64  `name:"min-mult" help:"Minimum traffic mult (0 means unlimited)." env:"MIN_MULT"`
}

// SharedCommandOptions AfterApply 在执行 Run 前运行
func (r *SharedCommandOptions) AfterApply() error {
	return nil
}

type compiledSharedConfig struct {
	SharedDirectCommandOptions
	NameRe     *regexp.Regexp
	RegionMap  map[string]string
	FilterSpec *filter.FilterSpec
}

// parseRegionMapEntries 解析地区映射
func parseRegionMapEntries(entries []string) (map[string]string, error) {
	result := make(map[string]string, len(entries))
	for _, entry := range entries {
		parts := strings.SplitN(entry, ":", 2)
		raw := strings.TrimSpace(parts[0])
		name := raw

		if len(parts) == 2 {
			name = strings.TrimSpace(parts[1])
		}
		if name == "" {
			name = raw
		}
		result[raw] = name
	}
	return result, nil
}

func (r *SharedCommandOptions) compile() (*compiledSharedConfig, error) {
	// 节点名正则
	pattern := model.DefaultNamePattern
	if r.NamePattern != "" {
		pattern = r.NamePattern
	}
	nameRe, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid name pattern: %w", err)
	}

	// 地区映射
	regionMap, err := parseRegionMapEntries(r.RegionMap)
	if err != nil {
		return nil, err
	}

	// 节点名排除正则
	var excludeReS []*regexp.Regexp
	for _, p := range r.ExcludePattern {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %w", err)
		}
		excludeReS = append(excludeReS, re)
	}

	return &compiledSharedConfig{
		SharedDirectCommandOptions: r.SharedDirectCommandOptions,
		NameRe:                     nameRe,
		RegionMap:                  regionMap,
		FilterSpec: &filter.FilterSpec{
			IncludeRegions: util.SliceToBoolMap(r.IncludeRegion),
			ExcludeRegions: util.SliceToBoolMap(r.ExcludeRegion),
			ExcludeType:    util.SliceToBoolMap(r.ExcludeType),
			ExcludeReS:     excludeReS,
			MaxMult:        r.MaxMult,
			MinMult:        r.MinMult,
		},
	}, nil
}
