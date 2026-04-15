package main

import (
	"fmt"
	"strings"

	"github.com/wuruimiao/clash-forge/internal/util"
	"github.com/wuruimiao/clash-forge/pkg/model"
)

type ListCmd struct {
	SharedCommandOptions `embed:""`
	Format               string `name:"format" default:"table" enum:"table,json" help:"List output format (${enum})." env:"LIST_FORMAT"`
}

func (r *ListCmd) Run() error {
	compiled, err := r.SharedCommandOptions.compile()
	if err != nil {
		return err
	}
	result, err := compiled.executePipeline()
	if err != nil {
		return err
	}
	switch r.Format {
	case "table":
		printTable(result.filteredNodes.NormalNodes)
		return nil
	case "json":
		return printJSON(result.filteredNodes.NormalNodes)
	default:
		return fmt.Errorf("unsupported list format: %s", r.Format)
	}
}

// printTable 以 ASCII 表格形式输出节点列表（调试用）。
func printTable(nodes []*model.Node) {
	// 定义表头宽度
	wIdx, wName, wReg, wMult := 6, 50, 30, 8
	// 打印表头
	fmt.Printf("%s %s %s %s",
		util.Pad("INDEX", wIdx),
		util.Pad("NAME", wName),
		util.Pad("REGION", wReg),
		util.Pad("MULT", wMult),
	)
	fmt.Println()
	fmt.Println(strings.Repeat("-", 105))

	for i, n := range nodes {
		fmt.Printf("%s %s %s %s\n",
			util.Pad(fmt.Sprintf("%d", i+1), wIdx),
			util.Pad(n.Name, wName),
			util.Pad(n.Region, wReg),
			util.Pad(n.Mult, wMult),
		)
	}
	fmt.Printf("\nTotal: %d nodes\n", len(nodes))
}

// printJSON 以 JSON 格式输出节点列表（调试用）。
func printJSON(nodes []*model.Node) error {
	type out struct {
		Name   string  `json:"name"`
		Region string  `json:"region"`
		Mult   float64 `json:"mult"`
		Type   string  `json:"type"`
		Server string  `json:"server"`
		Port   int     `json:"port"`
	}

	list := make([]out, 0, len(nodes))
	for _, n := range nodes {
		list = append(list, out{
			Name: n.Name, Region: n.Region,
			Mult: n.Mult, Type: n.Type, Server: n.Server, Port: n.Port,
		})
	}
	return util.JsonV(list)
}
