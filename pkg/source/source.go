// Package source 负责从远程 URL 或本地文件获取代理配置原始数据。
package source

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// FetchOne 获取单个配置源的原始内容。
// 自动识别输入类型：以 http:// 或 https:// 开头则通过 HTTP 获取，否则作为本地文件路径读取。
func FetchOne(input string) ([]byte, error) {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return fetchURL(input)
	}
	return os.ReadFile(input)
}

// fetchURL 通过 HTTP GET 请求获取远程配置文件。
// 使用 30 秒超时和浏览器伪装，避免被目标服务器拒绝。
func fetchURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 浏览器伪装
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response %s: %w", url, err)
	}
	return data, nil
}
