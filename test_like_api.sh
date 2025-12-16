#!/bin/bash

# 文章点赞功能API测试脚本
# 请确保服务器已启动并运行在 localhost:8888

BASE_URL="http://localhost:3100/backend"
API_BASE="$BASE_URL/like"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始测试文章点赞功能API...${NC}"
echo "=========================================="

# 测试用的帖子ID（请根据实际情况修改）
POST_ID=1
USER_TOKEN="your_jwt_token_here"  # 请替换为实际的JWT token

# 检查服务器是否运行
echo -e "${YELLOW}1. 检查服务器状态...${NC}"
curl -s -o /dev/null -w "%{http_code}" "$BASE_URL" > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 服务器运行正常${NC}"
else
    echo -e "${RED}✗ 服务器未运行，请先启动服务器${NC}"
    exit 1
fi

echo ""

# 测试1: 获取帖子点赞数
echo -e "${YELLOW}2. 测试获取帖子点赞数...${NC}"
echo "GET $API_BASE/getPostLikeCount/$POST_ID"
response=$(curl -s -w "\n%{http_code}" "$API_BASE/getPostLikeCount/$POST_ID")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

echo "HTTP状态码: $http_code"
echo "响应内容: $body"
if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ 获取帖子点赞数成功${NC}"
else
    echo -e "${RED}✗ 获取帖子点赞数失败${NC}"
fi

echo ""

# 测试2: 检查用户点赞状态（需要认证）
echo -e "${YELLOW}3. 测试检查用户点赞状态...${NC}"
echo "GET $API_BASE/checkUserLiked/$POST_ID"
if [ "$USER_TOKEN" != "your_jwt_token_here" ]; then
    response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $USER_TOKEN" "$API_BASE/checkUserLiked/$POST_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "HTTP状态码: $http_code"
    echo "响应内容: $body"
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 检查用户点赞状态成功${NC}"
    else
        echo -e "${RED}✗ 检查用户点赞状态失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 跳过测试（需要设置有效的JWT token）${NC}"
fi

echo ""

# 测试3: 获取帖子点赞列表
echo -e "${YELLOW}4. 测试获取帖子点赞列表...${NC}"
echo "GET $API_BASE/getPostLikes/$POST_ID?page=1&page_size=10"
response=$(curl -s -w "\n%{http_code}" "$API_BASE/getPostLikes/$POST_ID?page=1&page_size=10")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

echo "HTTP状态码: $http_code"
echo "响应内容: $body"
if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ 获取帖子点赞列表成功${NC}"
else
    echo -e "${RED}✗ 获取帖子点赞列表失败${NC}"
fi

echo ""

# 测试4: 点赞帖子（需要认证）
echo -e "${YELLOW}5. 测试点赞帖子...${NC}"
echo "POST $API_BASE/likePost/$POST_ID"
if [ "$USER_TOKEN" != "your_jwt_token_here" ]; then
    response=$(curl -s -w "\n%{http_code}" -X POST -H "Authorization: Bearer $USER_TOKEN" "$API_BASE/likePost/$POST_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "HTTP状态码: $http_code"
    echo "响应内容: $body"
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 点赞帖子成功${NC}"
    else
        echo -e "${RED}✗ 点赞帖子失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 跳过测试（需要设置有效的JWT token）${NC}"
fi

echo ""

# 测试5: 取消点赞帖子（需要认证）
echo -e "${YELLOW}6. 测试取消点赞帖子...${NC}"
echo "DELETE $API_BASE/unlikePost/$POST_ID"
if [ "$USER_TOKEN" != "your_jwt_token_here" ]; then
    response=$(curl -s -w "\n%{http_code}" -X DELETE -H "Authorization: Bearer $USER_TOKEN" "$API_BASE/unlikePost/$POST_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "HTTP状态码: $http_code"
    echo "响应内容: $body"
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 取消点赞帖子成功${NC}"
    else
        echo -e "${RED}✗ 取消点赞帖子失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 跳过测试（需要设置有效的JWT token）${NC}"
fi

echo ""

# 测试6: 获取用户点赞的帖子列表（需要认证）
echo -e "${YELLOW}7. 测试获取用户点赞的帖子列表...${NC}"
echo "GET $API_BASE/getUserLikedPosts?page=1&page_size=10"
if [ "$USER_TOKEN" != "your_jwt_token_here" ]; then
    response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $USER_TOKEN" "$API_BASE/getUserLikedPosts?page=1&page_size=10")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "HTTP状态码: $http_code"
    echo "响应内容: $body"
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 获取用户点赞的帖子列表成功${NC}"
    else
        echo -e "${RED}✗ 获取用户点赞的帖子列表失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 跳过测试（需要设置有效的JWT token）${NC}"
fi

echo ""
echo "=========================================="
echo -e "${YELLOW}测试完成！${NC}"
echo ""
echo -e "${YELLOW}使用说明：${NC}"
echo "1. 请确保服务器已启动"
echo "2. 修改脚本中的 POST_ID 为实际的帖子ID"
echo "3. 修改脚本中的 USER_TOKEN 为有效的JWT token"
echo "4. 运行脚本: bash test_like_api.sh"
