package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// URL映射结构
type URLMapping struct {
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	AccessCount int       `json:"access_count"`
}

// 存储结构
type URLStore struct {
	mappings map[string]*URLMapping
	mutex    sync.RWMutex
}

// 创建短链接请求
type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

// 创建短链接响应
type CreateShortURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// 统计信息响应
type StatsResponse struct {
	TotalURLs   int           `json:"total_urls"`
	TotalClicks int           `json:"total_clicks"`
	URLMappings []*URLMapping `json:"url_mappings"`
}

// 全局存储实例
var urlStore = &URLStore{
	mappings: make(map[string]*URLMapping),
}

// 生成短代码
func generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}

// 验证URL格式
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// 创建短链接
func createShortURL(c *gin.Context) {
	var req CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	if !isValidURL(req.OriginalURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供有效的URL（必须以http://或https://开头）"})
		return
	}

	// 生成短代码
	shortCode := generateShortCode()

	// 确保短代码唯一
	urlStore.mutex.Lock()
	for {
		if _, exists := urlStore.mappings[shortCode]; !exists {
			break
		}
		shortCode = generateShortCode()
	}

	// 创建映射
	mapping := &URLMapping{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		AccessCount: 0,
	}

	urlStore.mappings[shortCode] = mapping
	urlStore.mutex.Unlock()

	// 构建短URL
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080"
	}
	shortURL := fmt.Sprintf("http://%s/r/%s", host, shortCode)

	response := CreateShortURLResponse{
		ShortCode:   shortCode,
		ShortURL:    shortURL,
		OriginalURL: req.OriginalURL,
	}

	c.JSON(http.StatusCreated, response)
}

// 重定向到原始URL
func redirectToOriginal(c *gin.Context) {
	shortCode := c.Param("code")

	urlStore.mutex.RLock()
	mapping, exists := urlStore.mappings[shortCode]
	urlStore.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "短链接不存在"})
		return
	}

	// 增加访问计数
	urlStore.mutex.Lock()
	mapping.AccessCount++
	urlStore.mutex.Unlock()

	c.Redirect(http.StatusMovedPermanently, mapping.OriginalURL)
}

// 获取统计信息
func getStats(c *gin.Context) {
	urlStore.mutex.RLock()
	defer urlStore.mutex.RUnlock()

	var totalClicks int
	var mappings []*URLMapping

	for _, mapping := range urlStore.mappings {
		totalClicks += mapping.AccessCount
		mappings = append(mappings, mapping)
	}

	response := StatsResponse{
		TotalURLs:   len(urlStore.mappings),
		TotalClicks: totalClicks,
		URLMappings: mappings,
	}

	c.JSON(http.StatusOK, response)
}

// 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "URL Shortener",
	})
}

// 设置路由
func setupRoutes() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// API路由
	api := r.Group("/api/v1")
	{
		api.POST("/shorten", createShortURL)
		api.GET("/stats", getStats)
	}

	// 重定向路由
	r.GET("/r/:code", redirectToOriginal)

	// 健康检查
	r.GET("/health", healthCheck)

	// 根路径显示服务信息
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "URL Shortener",
			"version":     "1.0.0",
			"description": "一个简单的URL短链接服务",
			"endpoints": gin.H{
				"create_short_url": "POST /api/v1/shorten",
				"get_stats":        "GET /api/v1/stats",
				"redirect":         "GET /r/:code",
				"health_check":     "GET /health",
			},
		})
	})

	return r
}

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 设置路由
	r := setupRoutes()

	// 启动服务器
	port := ":8080"
	log.Printf("URL短链接服务启动在端口 %s", port)
	log.Printf("访问 http://localhost%s 查看API文档", port)

	if err := r.Run(port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
