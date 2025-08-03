# URL短链接微服务

一个使用Go语言和Gin框架构建的简单URL短链接微服务。

## 功能特性

- ✅ 创建短链接
- ✅ 重定向到原始URL
- ✅ 访问统计
- ✅ RESTful API
- ✅ 健康检查
- ✅ CORS支持
- ✅ 并发安全

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 3. 测试服务

访问 `http://localhost:8080` 查看API文档。

## API文档

### 1. 创建短链接

**请求:**
```http
POST /api/v1/shorten
Content-Type: application/json

{
    "original_url": "https://www.example.com/very/long/url"
}
```

**响应:**
```json
{
    "short_code": "abc12345",
    "short_url": "http://localhost:8080/r/abc12345",
    "original_url": "https://www.example.com/very/long/url"
}
```

### 2. 重定向到原始URL

**请求:**
```http
GET /r/{short_code}
```

**响应:**
- 302重定向到原始URL

### 3. 获取统计信息

**请求:**
```http
GET /api/v1/stats
```

**响应:**
```json
{
    "total_urls": 5,
    "total_clicks": 25,
    "url_mappings": [
        {
            "original_url": "https://www.example.com",
            "short_code": "abc12345",
            "created_at": "2024-01-01T12:00:00Z",
            "access_count": 10
        }
    ]
}
```

### 4. 健康检查

**请求:**
```http
GET /health
```

**响应:**
```json
{
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "service": "URL Shortener"
}
```

## 使用示例

### 使用curl创建短链接

```bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://www.google.com"}'
```

### 使用curl获取统计信息

```bash
curl http://localhost:8080/api/v1/stats
```

### 在浏览器中测试重定向

访问生成的短链接，例如：`http://localhost:8080/r/abc12345`

## 项目结构

```
URL_shortener/
├── main.go          # 主程序文件
├── go.mod           # Go模块文件
└── README.md        # 项目说明
```

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **UUID生成**: google/uuid
- **存储**: 内存存储（可扩展为数据库）

## 扩展建议

1. **持久化存储**: 集成Redis或PostgreSQL
2. **自定义短代码**: 允许用户指定短代码
3. **过期时间**: 为短链接设置过期时间
4. **用户认证**: 添加用户系统
5. **API限流**: 防止滥用
6. **监控和日志**: 添加详细的监控和日志
7. **Docker化**: 创建Docker镜像

## 许可证

MIT License 