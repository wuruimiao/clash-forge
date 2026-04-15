package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFileAtomic 通过"先写临时文件再重命名"的方式安全写入文件。
// 这样即使写入过程中断，也不会留下损坏的目标文件。
func WriteFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	// 创建目录
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dir for region yaml: %w", err)
		}
	}

	// 临时文件
	tmp, err := os.CreateTemp(dir, ".clash-forge-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	// 写入临时文件
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
