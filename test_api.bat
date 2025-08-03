@echo off
chcp 65001 >nul
echo 🧪 开始测试URL短链接服务...
echo ==================================

set BASE_URL=http://localhost:8080

echo 1. 测试健康检查...
curl -s "%BASE_URL%/health"
echo.
echo.

echo 2. 获取服务信息...
curl -s "%BASE_URL%/"
echo.
echo.

echo 3. 创建短链接...
curl -s -X POST "%BASE_URL%/api/v1/shorten" -H "Content-Type: application/json" -d "{\"original_url\": \"https://www.google.com\"}"
echo.
echo.

echo 4. 获取统计信息...
curl -s "%BASE_URL%/api/v1/stats"
echo.
echo.

echo ✅ 测试完成！
echo.
echo 📝 使用说明:
echo - 服务运行在: %BASE_URL%
echo - 创建短链接: POST %BASE_URL%/api/v1/shorten
echo - 查看统计: GET %BASE_URL%/api/v1/stats
echo - 健康检查: GET %BASE_URL%/health
echo - 重定向: GET %BASE_URL%/r/{short_code}
echo.
pause 