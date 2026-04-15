package main

import (
	"fmt"
	"strings"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

type PreviewCmd struct {
	SharedCommandOptions `embed:""`
}

type VersionCmd struct{}

func (r *PreviewCmd) Run() error {
	compiled, err := r.SharedCommandOptions.compile()
	if err != nil {
		return err
	}
	result, err := compiled.executePipeline()
	if err != nil {
		return err
	}
	printPreview(result.allNodes)
	return nil
}

// printPreview 展示所有节点的解析结果和分组摘要（含 info 节点）。
func printPreview(nodes []*model.Node) {
	// 定义表头宽度
	wIdx, wName, wReg, wMult := 6, 50, 30, 8

	// 打印表头
	fmt.Printf("%s %s %s %s %s",
		util.Pad("INDEX", wIdx),
		util.Pad("NAME", wName),
		util.Pad("REGION", wReg),
		util.Pad("MULT", wMult),
		"INFO")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 105))

	var proxyCount, infoCount int
	regionCounts := make(map[string]int)
	multCounts := make(map[float64]int)

	for i, n := range nodes {
		if n.IsInfo {
			infoCount++
			fmt.Printf("%s %s %s %s %s\n",
				util.Pad(fmt.Sprintf("%d", i+1), wIdx),
				util.Pad(n.Name, wName),
				util.Pad("-", wReg),
				util.Pad("-", wMult),
				"Yes")
			// fmt.Printf("%-12d %-45s %-10s %-8s %s\n", i+1, n.Name, "-", "-", "Yes")
		} else {
			proxyCount++
			regionCounts[n.Region]++
			multCounts[n.Mult]++
			fmt.Printf("%s %s %s %s %s\n",
				util.Pad(fmt.Sprintf("%d", i+1), wIdx),
				util.Pad(n.Name, wName),
				util.Pad(n.Region, wReg),
				util.Pad(fmt.Sprintf("%0.1f", n.Mult), wMult),
				"No")
			// fmt.Printf("%-12d %-45s %-10s %-8.1f\n", i+1, n.Name, n.Region, n.Mult)
		}
	}

	fmt.Printf("\nTotal: %d nodes (%d proxy, %d info)\n", len(nodes), proxyCount, infoCount)

	if len(regionCounts) > 0 {
		fmt.Printf("\nRegions (%d):\n", len(regionCounts))
		for region, count := range regionCounts {
			fmt.Printf("  %-10s %d nodes\n", region, count)
		}
	}

	if len(multCounts) > 0 {
		fmt.Printf("\nMults (%d):\n", len(multCounts))
		for mult, count := range multCounts {
			fmt.Printf("  %-10.1f %d nodes\n", mult, count)
		}
	}
}

func (r *VersionCmd) Run() error {
	fmt.Println(version)
	return nil
}
