package tsplay_core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatClient 结构体封装了与聊天API的交互逻辑
type ChatClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewChatClient 创建一个新的ChatClient实例
func NewChatClient(baseURL string) *ChatClient {
	return &ChatClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// Message 结构体表示API消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 结构体表示API请求体
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Options  Options   `json:"options"`
}

// Options 结构体表示API请求的选项
type Options struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ChatResponse 结构体表示API响应体
type ChatResponse struct {
	Message Message `json:"message"`
}

// SendChatRequest 发送聊天请求并返回响应
func (c *ChatClient) SendChatRequest(model string, messages []Message, stream bool, options Options) (*ChatResponse, error) {
	// 构建请求体
	requestBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   stream,
		Options:  options,
	}

	// 将请求体编码为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", c.BaseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// 解析响应体
	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return &response, nil
}

// 配置合适的promt，并封装出合适的功能

//

// 案例代码
//package main
//
//import (
//"fmt"
//"your_project/chatclient" // 替换为你的模块路径
//)
//
//func main() {
//	// 创建ChatClient实例
//	client := chatclient.NewChatClient("http://127.0.0.1:11436")
//
//	// 定义请求消息
//	messages := []chatclient.Message{
//		{
//			Role:    "user",
//			Content: "初中生必备古诗词",
//		},
//	}
//
//	// 定义请求选项
//	options := chatclient.Options{
//		Temperature: 0.7,
//		MaxTokens:   100,
//	}
//
//	// 发送请求
//	response, err := client.SendChatRequest("deepseek-v3:671b(sdu)", messages, false, options)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// 打印响应
//	fmt.Println("Response:", response.Message.Content)
//}
