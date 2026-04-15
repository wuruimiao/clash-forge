package util

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
)

// MatchesAny 检查字符串是否匹配任意一个正则模式。
func MatchesAny(patterns []*regexp.Regexp, s string) bool {
	for _, re := range patterns {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

// ToInt 将 YAML 解析出的 any 类型安全转换为 int。
// YAML 数字可能被解析为 int、float64 或 int64，需要逐一处理。
func ToInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	case int64:
		return int(n)
	default:
		return 0
	}
}

func ToString(v any) string {
	switch val := v.(type) {
	case float32, float64:
		return fmt.Sprintf("%g", val) // 浮点数使用 %g
	default:
		return fmt.Sprintf("%v", val) // 字符串、整数等使用 %v
	}
}

// SortedKeys 提取 map 的键并按字典序从小到大排序，确保输出顺序稳定。
func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// StrVal 从 map 中安全提取字符串值，键不存在或类型不匹配时返回空字符串。
func StrVal(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

// StrValOr 从 map 中提取字符串值，值为空时返回 fallback 默认值。
func StrValOr(m map[string]any, key, fallback string) string {
	if v := StrVal(m, key); v != "" {
		return v
	}
	return fallback
}

// Map2Struct 泛型函数：将 map 提取到指定的结构体 T 中
func Map2Struct[T any](raw map[string]any) (T, error) {
	var result T

	// 1. 将 map 序列化为二进制 JSON
	// 这是最简单且兼容性最好的方式，能自动处理 map 到 struct 标签的映射
	bytes, err := json.Marshal(raw)
	if err != nil {
		return result, err
	}

	// 2. 反序列化到目标结构体类型
	err = json.Unmarshal(bytes, &result)
	return result, err
}

func Unique[T comparable](input []T) []T {
	result := make([]T, 0, len(input))
	seen := make(map[T]struct{})

	for _, val := range input {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}
func SliceToBoolMap(values []string) map[string]bool {
	if len(values) == 0 {
		return nil
	}
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func JsonV(value any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
