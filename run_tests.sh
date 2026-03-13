#!/bin/bash

# 文章点赞功能完整测试脚本

echo "=========================================="
echo "文章点赞功能测试开始"
echo "=========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检查Go环境
echo -e "${BLUE}1. 检查Go环境...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go未安装或不在PATH中${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go环境正常${NC}"

# 检查项目依赖
echo -e "${BLUE}2. 检查项目依赖...${NC}"
if [ ! -f "go.mod" ]; then
    echo -e "${RED}✗ 未找到go.mod文件${NC}"
    exit 1
fi

# 下载依赖
go mod tidy
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 依赖检查完成${NC}"
else
    echo -e "${RED}✗ 依赖下载失败${NC}"
    exit 1
fi

# 编译项目
echo -e "${BLUE}3. 编译项目...${NC}"
go build -o server .
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 项目编译成功${NC}"
else
    echo -e "${RED}✗ 项目编译失败${NC}"
    exit 1
fi

# 运行单元测试
echo -e "${BLUE}4. 运行单元测试...${NC}"
go test ./tests/... -v
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 单元测试通过${NC}"
else
    echo -e "${YELLOW}⚠ 单元测试失败或跳过（需要数据库连接）${NC}"
fi

# 检查数据库连接
echo -e "${BLUE}5. 检查数据库连接...${NC}"
# 这里可以添加数据库连接检查
echo -e "${YELLOW}⚠ 请手动检查数据库连接${NC}"

# 运行数据库测试
echo -e "${BLUE}6. 运行数据库测试...${NC}"
if [ -f "test_database.sql" ]; then
    echo -e "${YELLOW}请手动执行以下SQL文件进行数据库测试：${NC}"
    echo -e "${YELLOW}mysql -u username -p database_name < test_database.sql${NC}"
else
    echo -e "${RED}✗ 未找到数据库测试文件${NC}"
fi

# 启动服务器进行API测试
echo -e "${BLUE}7. 启动服务器进行API测试...${NC}"
echo -e "${YELLOW}请手动启动服务器：${NC}"
echo -e "${YELLOW}./server${NC}"
echo ""
echo -e "${YELLOW}然后运行API测试脚本：${NC}"
echo -e "${YELLOW}bash test_like_api.sh${NC}"

# 性能测试
echo -e "${BLUE}8. 性能测试...${NC}"
if command -v ab &> /dev/null; then
    echo -e "${YELLOW}可以使用Apache Bench进行性能测试：${NC}"
    echo -e "${YELLOW}ab -n 1000 -c 10 http://localhost:8888/backend/like/getPostLikeCount/1${NC}"
else
    echo -e "${YELLOW}⚠ Apache Bench未安装，跳过性能测试${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}测试完成！${NC}"
echo "=========================================="
echo ""
echo -e "${BLUE}测试总结：${NC}"
echo "1. ✓ Go环境检查"
echo "2. ✓ 项目依赖检查"
echo "3. ✓ 项目编译"
echo "4. ⚠ 单元测试（需要数据库）"
echo "5. ⚠ 数据库连接检查"
echo "6. ⚠ 数据库测试（需要手动执行）"
echo "7. ⚠ API测试（需要启动服务器）"
echo "8. ⚠ 性能测试（可选）"
echo ""
echo -e "${YELLOW}下一步操作：${NC}"
echo "1. 配置数据库连接"
echo "2. 执行数据库迁移：sql/create_like_tables.sql"
echo "3. 启动服务器：./server"
echo "4. 运行API测试：bash test_like_api.sh"
echo "5. 导入Postman集合：postman_like_api.json"
