package baiduyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// 百度开放平台OAuth2授权地址
	authURL = "https://openapi.baidu.com/oauth/2.0/token"
	// 百度云盘OpenAPI地址
	apiURL = "https://pan.baidu.com/rest/2.0/xpan/file"
)

// Client 百度云盘客户端
type Client struct {
	clientID     string
	clientSecret string
	redirectURI  string
	accessToken  string
	refreshToken string
	httpClient   *http.Client
}

// NewClient 创建百度云盘客户端
func NewClient(clientID, clientSecret, redirectURI string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		httpClient:   &http.Client{},
	}
}

// NewClientWithToken 使用已有AccessToken创建客户端
func NewClientWithToken(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

// GetAuthURL 获取授权页面URL
func (c *Client) GetAuthURL() string {
	return fmt.Sprintf(
		"https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=basic,netdisk",
		c.clientID,
		url.QueryEscape(c.redirectURI),
	)
}

// ExchangeToken 用授权码交换AccessToken
func (c *Client) ExchangeToken(code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("redirect_uri", c.redirectURI)

	resp, err := c.httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf("授权失败: %s - %s", result.Error, result.ErrorDesc)
	}

	c.accessToken = result.AccessToken
	c.refreshToken = result.RefreshToken
	return nil
}

// RefreshToken 刷新AccessToken
func (c *Client) RefreshToken() error {
	if c.refreshToken == "" {
		return errors.New("refresh token is empty")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.refreshToken)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	resp, err := c.httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf("刷新token失败: %s - %s", result.Error, result.ErrorDesc)
	}

	c.accessToken = result.AccessToken
	c.refreshToken = result.RefreshToken
	return nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://pan.baidu.com/rest/2.0/xpan/nas?method=uinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if errCode, ok := result["errno"].(float64); ok && errCode != 0 {
		return nil, fmt.Errorf("获取用户信息失败，错误码: %.0f", errCode)
	}

	return result, nil
}

// AccessToken 获取当前AccessToken
func (c *Client) AccessToken() string {
	return c.accessToken
}
