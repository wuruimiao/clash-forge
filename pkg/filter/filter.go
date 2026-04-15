// Package filter 根据用户指定的条件对代理节点进行过滤。
// 支持按地区（包含/排除）、倍率范围、正则模式等多维度筛选。
package filter

import (
	"log/slog"
	"regexp"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

// FilterSpec 定义节点过滤条件。
type FilterSpec struct {
	IncludeRegions map[string]bool  // 要包含的地区名
	ExcludeRegions map[string]bool  // 要排除的地区名
	ExcludeType    map[string]bool  // 要排除的节点类型（如 "vmess", "trojan"）
	MaxMult        float64          // 最大倍率上限（0 表示不限制）
	MinMult        float64          // 最小倍率下限（0 表示不限制）
	ExcludeReS     []*regexp.Regexp // 正则排除模式，匹配节点名称
}

// Apply 根据 FilterSpec 过滤节点列表。
// 过滤顺序：跳过信息节点 → 地区白名单 → 地区黑名单 → 倍率下限 → 倍率上限 → 正则排除。
func (f *FilterSpec) Apply(nodes []*model.Node) []*model.Node {
	var result []*model.Node
	for _, n := range nodes {
		// 排除指定类型的节点
		if f.ExcludeType[n.Type] {
			slog.Warn("skipping proxy: excluded type", "name", n.Name)
			continue
		}
		// 如果设置了白名单，不在白名单中的节点跳过
		if len(f.IncludeRegions) > 0 && !f.IncludeRegions[n.Region] {
			slog.Warn("skipping proxy: excluded region", "name", n.Name)
			continue
		}
		// 在黑名单中的节点跳过
		if f.ExcludeRegions[n.Region] {
			slog.Warn("skipping proxy: excluded region", "name", n.Name)
			continue
		}
		// 倍率低于下限的跳过
		if f.MinMult > 0 && n.Mult < f.MinMult {
			slog.Warn("skipping proxy: excluded mult", "name", n.Name)
			continue
		}
		// 倍率超过上限的跳过
		if f.MaxMult > 0 && n.Mult > f.MaxMult {
			slog.Warn("skipping proxy: excluded mult", "name", n.Name)
			continue
		}
		// 节点名匹配任一排除正则的跳过
		if util.MatchesAny(f.ExcludeReS, n.Name) {
			slog.Warn("skipping proxy: excluded regex", "name", n.Name)
			continue
		}
		result = append(result, n)
	}
	return result
}
