// Package grouper 根据过滤后的节点列表自动生成代理分组。
// 生成三类分组：全部节点（手动选择）、按地区分组（自动测速）、按倍率分组（自动测速）。
package grouper

import (
	"fmt"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

// Group 根据过滤后的节点列表创建代理分组。
func Group(nodes []*model.Node) *model.NodeGroups {
	if len(nodes) == 0 {
		return nil
	}

	groups := make(map[model.NodeGroupType][]*model.NodeGroup)

	// 按地区分组
	regionGroups := groupBy(nodes, model.Region, func(n *model.Node) string {
		return n.Region
	})
	groups[model.Region] = append(groups[model.Region], regionGroups...)

	// 按倍率分组
	multGroups := groupBy(nodes, model.Mult, func(n *model.Node) string {
		// formatMult 将倍率数值格式化为中文显示名。
		// 整数显示为 "倍率1"，小数显示为 "倍率0.1"。
		if n.Mult == float64(int(n.Mult)) {
			return fmt.Sprintf("倍率%d", int(n.Mult))
		}
		return fmt.Sprintf("倍率%.1f", n.Mult)
	})
	groups[model.Mult] = append(groups[model.Mult], multGroups...)

	return &model.NodeGroups{
		Groups: groups,
	}
}

// groupby 根据getKey+category，对nodes分组，一个 url-test 组，自动选择延迟最低的节点
func groupBy(nodes []*model.Node, groupType model.NodeGroupType, getKey func(*model.Node) string) []*model.NodeGroup {
	nodeMap := map[string][]*model.Node{}
	for _, n := range nodes {
		key := getKey(n)
		if key != "" {
			nodeMap[key] = append(nodeMap[key], n)
		}
	}

	var res []*model.NodeGroup
	for _, key := range util.SortedKeys(nodeMap) {
		res = append(res, &model.NodeGroup{
			Name:  key,
			Type:  groupType,
			Nodes: nodeMap[key],
		})
	}
	return res
}
