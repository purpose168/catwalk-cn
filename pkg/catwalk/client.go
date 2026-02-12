package catwalk

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	xetag "github.com/purpose168/charm-experimental-packages-cn/etag"
)

const defaultURL = "http://localhost:8080"

// Client 表示catwalk服务的客户端
type Client struct {
	baseURL    string       // 服务基础URL
	httpClient *http.Client // HTTP客户端实例
}

// New 创建一个新的客户端实例
// 使用CATWALK_URL环境变量，如果未设置则回退到localhost:8080
func New() *Client {
	return NewWithURL(cmp.Or(os.Getenv("CATWALK_URL"), defaultURL))
}

// NewWithURL 使用指定的URL创建一个新的客户端
func NewWithURL(url string) *Client {
	return &Client{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ErrNotModified 当给定的ETag与服务器匹配时发生，表示不需要更新
var ErrNotModified = fmt.Errorf("未修改")

// Etag 返回给定数据的ETag
func Etag(data []byte) string { return xetag.Of(data) }

// GetProviders 从服务中检索所有可用的AI提供商
func (c *Client) GetProviders(ctx context.Context, etag string) ([]Provider, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v2/providers", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %w", err)
	}
	xetag.Request(req, etag)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusNotModified {
		return nil, ErrNotModified
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("意外的状态码: %d", resp.StatusCode)
	}

	var providers []Provider
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, fmt.Errorf("响应解码失败: %w", err)
	}

	return providers, nil
}
