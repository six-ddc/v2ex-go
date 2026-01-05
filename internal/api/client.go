package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

const (
	BaseURL   = "https://www.v2ex.com"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Client V2EX HTTP 客户端
type Client struct {
	httpClient *resty.Client
	baseURL    string
	cookies    []*http.Cookie
}

// NewClient 创建新的 HTTP 客户端
func NewClient() *Client {
	client := resty.New().
		SetBaseURL(BaseURL).
		SetTimeout(30 * time.Second).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	return &Client{
		httpClient: client,
		baseURL:    BaseURL,
	}
}

// SetCookies 设置登录 Cookie
func (c *Client) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
	c.httpClient.SetCookies(cookies)
}

// Get 发送 GET 请求并返回解析后的 HTML 文档
func (c *Client) Get(path string) (*goquery.Document, error) {
	resp, err := c.httpClient.R().Get(path)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// GetRaw 发送 GET 请求并返回原始响应
func (c *Client) GetRaw(path string) (*resty.Response, error) {
	return c.httpClient.R().Get(path)
}

// Post 发送 POST 请求并返回解析后的 HTML 文档
func (c *Client) Post(path string, data map[string]string) (*goquery.Document, error) {
	resp, err := c.httpClient.R().
		SetFormData(data).
		Post(path)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// BuildURL 构建完整 URL
func (c *Client) BuildURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return c.baseURL + path
}
