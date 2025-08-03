#!/bin/bash

# URL短链接服务API测试脚本

BASE_URL="http://localhost:8080"

echo "🧪 开始测试URL短链接服务..."
echo "=================================="

# 测试健康检查
echo "1. 测试健康检查..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 测试服务信息
echo "2. 获取服务信息..."
curl -s "$BASE_URL/" | jq .
echo ""

# 测试创建短链接
echo "3. 创建短链接..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/shorten" \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://www.google.com"}')

echo "$RESPONSE" | jq .

# 提取短代码
SHORT_CODE=$(echo "$RESPONSE" | jq -r '.short_code')
SHORT_URL=$(echo "$RESPONSE" | jq -r '.short_url')

echo ""
echo "生成的短链接: $SHORT_URL"
echo ""

# 测试获取统计信息
echo "4. 获取统计信息..."
curl -s "$BASE_URL/api/v1/stats" | jq .
echo ""

# 测试重定向（只检查状态码）
echo "5. 测试重定向（检查状态码）..."
REDIRECT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/r/$SHORT_CODE")
echo "重定向状态码: $REDIRECT_STATUS"
echo ""

echo "✅ 测试完成！"
echo ""
echo "📝 使用说明:"
echo "- 服务运行在: $BASE_URL"
echo "- 创建短链接: POST $BASE_URL/api/v1/shorten"
echo "- 查看统计: GET $BASE_URL/api/v1/stats"
echo "- 健康检查: GET $BASE_URL/health"
echo "- 重定向: GET $BASE_URL/r/{short_code}" 