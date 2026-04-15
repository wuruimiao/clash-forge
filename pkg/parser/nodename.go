package parser

import (
	"log/slog"
	"regexp"
	"strconv"

	"github.com/wuruimiao/clash-forge/pkg/model"
)

type NodeNameParser struct {
	re        *regexp.Regexp
	multIdx   int
	regionIdx int
	trie      *regionTrie
}

type NodeNameParserOptions struct {
	Re           *regexp.Regexp
	ExtraRegions map[string]string
}

func (n *NodeNameParserOptions) New() (*NodeNameParser, error) {
	trie := buildTrie(n.ExtraRegions)
	p := &NodeNameParser{
		re:        n.Re,
		multIdx:   -1,
		regionIdx: -1,
		trie:      trie,
	}
	for i, name := range n.Re.SubexpNames() {
		switch name {
		case "mult":
			p.multIdx = i
		case "region":
			p.regionIdx = i
		}
	}
	return p, nil
}

func (p *NodeNameParser) Parse(name string) model.NodeNameMeta {
	meta := model.NodeNameMeta{IsInfo: true}
	if name == "" {
		return meta
	}
	matches := p.re.FindStringSubmatch(name)
	if matches == nil {
		return meta
	}
	multStr := matches[p.multIdx]
	mult, err := strconv.ParseFloat(multStr, 64)
	if err != nil {
		slog.Warn("no mult as info", "name", name)
		return meta
	}
	meta.IsInfo = false
	meta.Mult = mult
	if p.regionIdx >= 0 && p.regionIdx < len(matches) {
		meta.RegionRaw = matches[p.regionIdx]
		_, meta.Region, _ = p.trie.longestMatch([]rune(meta.RegionRaw))
		if meta.Region == "" {
			meta.Region = meta.RegionRaw
		}
	}
	return meta
}
