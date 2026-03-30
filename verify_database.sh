#!/bin/bash

# 数据库结构验证脚本
# 用于验证comments表结构是否正确更新

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 数据库配置
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_NAME=${DB_NAME:-"server_fiber"}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-"123456"}

echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}数据库结构验证工具${NC}"
echo -e "${BLUE}==========================================${NC}"

# 检查数据库连接
echo -e "${YELLOW}1. 检查数据库连接...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SELECT 1;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 数据库连接成功${NC}"
else
    echo -e "${RED}✗ 数据库连接失败${NC}"
    exit 1
fi

# 验证comments表结构
echo -e "${YELLOW}2. 验证comments表结构...${NC}"

# 检查表是否存在
TABLE_EXISTS=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SHOW TABLES LIKE 'comments';" 2>/dev/null | wc -l)
if [ "$TABLE_EXISTS" -le 1 ]; then
    echo -e "${RED}✗ comments表不存在${NC}"
    exit 1
fi

echo -e "${GREEN}✓ comments表存在${NC}"

# 检查必需字段
echo -e "${YELLOW}3. 检查必需字段...${NC}"

REQUIRED_FIELDS=("id" "created_at" "updated_at" "deleted_at" "post_id" "parent_id" "content" "user_id" "to_user_id")
MISSING_FIELDS=()

for field in "${REQUIRED_FIELDS[@]}"; do
    FIELD_EXISTS=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; DESCRIBE comments;" 2>/dev/null | grep -c "$field")
    if [ "$FIELD_EXISTS" -eq 0 ]; then
        MISSING_FIELDS+=("$field")
    fi
done

if [ ${#MISSING_FIELDS[@]} -eq 0 ]; then
    echo -e "${GREEN}✓ 所有必需字段都存在${NC}"
else
    echo -e "${RED}✗ 缺少字段: ${MISSING_FIELDS[*]}${NC}"
fi

# 检查字段类型
echo -e "${YELLOW}4. 检查字段类型...${NC}"

# 检查post_id字段类型
POST_ID_TYPE=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; DESCRIBE comments;" 2>/dev/null | grep "post_id" | awk '{print $2}')
if [[ "$POST_ID_TYPE" == *"bigint"* ]]; then
    echo -e "${GREEN}✓ post_id字段类型正确: $POST_ID_TYPE${NC}"
else
    echo -e "${RED}✗ post_id字段类型错误: $POST_ID_TYPE${NC}"
fi

# 检查parent_id字段类型
PARENT_ID_TYPE=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; DESCRIBE comments;" 2>/dev/null | grep "parent_id" | awk '{print $2}')
if [[ "$PARENT_ID_TYPE" == *"bigint"* ]]; then
    echo -e "${GREEN}✓ parent_id字段类型正确: $PARENT_ID_TYPE${NC}"
else
    echo -e "${RED}✗ parent_id字段类型错误: $PARENT_ID_TYPE${NC}"
fi

# 检查user_id字段类型
USER_ID_TYPE=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; DESCRIBE comments;" 2>/dev/null | grep "user_id" | awk '{print $2}')
if [[ "$USER_ID_TYPE" == *"bigint"* ]]; then
    echo -e "${GREEN}✓ user_id字段类型正确: $USER_ID_TYPE${NC}"
else
    echo -e "${RED}✗ user_id字段类型错误: $USER_ID_TYPE${NC}"
fi

# 检查索引
echo -e "${YELLOW}5. 检查索引...${NC}"

REQUIRED_INDEXES=("idx_comments_post_id" "idx_comments_parent_id" "idx_comments_user_id" "idx_comments_to_user_id")
MISSING_INDEXES=()

for index in "${REQUIRED_INDEXES[@]}"; do
    INDEX_EXISTS=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SHOW INDEX FROM comments;" 2>/dev/null | grep -c "$index")
    if [ "$INDEX_EXISTS" -eq 0 ]; then
        MISSING_INDEXES+=("$index")
    fi
done

if [ ${#MISSING_INDEXES[@]} -eq 0 ]; then
    echo -e "${GREEN}✓ 所有必需索引都存在${NC}"
else
    echo -e "${RED}✗ 缺少索引: ${MISSING_INDEXES[*]}${NC}"
fi

# 显示完整表结构
echo -e "${YELLOW}6. 完整表结构：${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; DESCRIBE comments;"

echo -e "${YELLOW}7. 索引信息：${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SHOW INDEX FROM comments;"

# 检查数据完整性
echo -e "${YELLOW}8. 数据完整性检查...${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SELECT 
  COUNT(*) as total_comments,
  COUNT(DISTINCT post_id) as unique_posts,
  COUNT(DISTINCT user_id) as unique_users,
  COUNT(CASE WHEN parent_id > 0 THEN 1 END) as reply_comments
FROM comments;"

# 总结
echo -e "${BLUE}==========================================${NC}"
if [ ${#MISSING_FIELDS[@]} -eq 0 ] && [ ${#MISSING_INDEXES[@]} -eq 0 ]; then
    echo -e "${GREEN}✓ 数据库结构验证通过！${NC}"
    echo -e "${GREEN}comments表结构已正确更新${NC}"
else
    echo -e "${RED}✗ 数据库结构验证失败！${NC}"
    echo -e "${YELLOW}请检查上述错误并重新运行迁移脚本${NC}"
fi
echo -e "${BLUE}==========================================${NC}"
