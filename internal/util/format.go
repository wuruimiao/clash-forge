package util

import (
	"fmt"
	"unicode"
)

// 计算字符串在终端显示的实际宽度（中文算2，英文算1）
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.In(r, unicode.Punct) && r > 127 {
			width += 2 // 中文字符和中文标点占2位
		} else {
			width += 1 // 英文字符、数字占1位
		}
	}
	return width
}

// 格式化输出函数：s 为内容，totalWidth 为预留的总宽度
func Pad(v any, totalWidth int) string {
	s := ToString(v)
	curWidth := getDisplayWidth(s)
	if curWidth >= totalWidth {
		return s
	}
	// 补足空格：总宽度 - 实际占位宽度
	return s + fmt.Sprintf("%*s", totalWidth-curWidth, "")
}
