package parser

import "log/slog"

type regionTrie struct {
	children    map[rune]*regionTrie
	raw         string
	mappingName string
}

func buildTrie(entries map[string]string) *regionTrie {
	root := &regionTrie{}
	for raw, name := range entries {
		if raw == "" {
			slog.Warn("empty raw")
			continue
		}
		node := root
		for _, r := range raw {
			if node.children == nil {
				node.children = map[rune]*regionTrie{}
			}
			child, ok := node.children[r]
			if !ok {
				child = &regionTrie{}
				node.children[r] = child
			}
			node = child
		}
		node.raw = raw
		node.mappingName = name
	}
	return root
}

func (t *regionTrie) longestMatch(s []rune) (raw string, mappingName string, matchLen int) {
	node := t
	for i, r := range s {
		child, ok := node.children[r]
		if !ok {
			break
		}
		node = child
		if node.raw != "" {
			raw = node.raw
			mappingName = node.mappingName
			matchLen = i + 1
		}
	}
	return
}
