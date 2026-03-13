package tests

import (
	"testing"

	global "server/model"
	"server/model/app"
	appService "server/service/app"

	"github.com/stretchr/testify/assert"
)

// TestLikeServiceBasic 测试点赞服务基础功能
func TestLikeServiceBasic(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过服务测试")
	// 初始化测试数据库
	setupTestDB()

	// 创建测试用户和帖子
	user := createTestUser("like_test_user")
	post := createTestPost(user.ID, "测试帖子", "这是一个测试帖子的内容")

	likeService := appService.LikeServer

	t.Run("测试点赞帖子", func(t *testing.T) {
		err := likeService.LikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 验证点赞记录已创建
		var like app.Like
		err = global.DB.Where("post_id = ? AND user_id = ?", post.ID, user.ID).First(&like).Error
		assert.NoError(t, err)
		assert.Equal(t, post.ID, like.PostId)
		assert.Equal(t, user.ID, like.UserId)

		// 验证帖子点赞数已更新
		var updatedPost app.Post
		err = global.DB.Where("id = ?", post.ID).First(&updatedPost).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, updatedPost.LikeCount)
	})

	t.Run("测试重复点赞", func(t *testing.T) {
		err := likeService.LikePost(post.ID, user.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经点赞过这篇帖子了")
	})

	t.Run("测试取消点赞", func(t *testing.T) {
		err := likeService.UnlikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 验证点赞记录已删除
		var count int64
		err = global.DB.Model(&app.Like{}).Where("post_id = ? AND user_id = ?", post.ID, user.ID).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// 验证帖子点赞数已更新
		var updatedPost app.Post
		err = global.DB.Where("id = ?", post.ID).First(&updatedPost).Error
		assert.NoError(t, err)
		assert.Equal(t, 0, updatedPost.LikeCount)
	})

	t.Run("测试取消不存在的点赞", func(t *testing.T) {
		err := likeService.UnlikePost(post.ID, user.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "还没有点赞过这篇帖子")
	})

	t.Run("测试检查用户点赞状态", func(t *testing.T) {
		// 先点赞
		err := likeService.LikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 检查点赞状态
		liked, err := likeService.CheckUserLiked(post.ID, user.ID)
		assert.NoError(t, err)
		assert.True(t, liked)

		// 取消点赞
		err = likeService.UnlikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 再次检查点赞状态
		liked, err = likeService.CheckUserLiked(post.ID, user.ID)
		assert.NoError(t, err)
		assert.False(t, liked)
	})

	t.Run("测试获取帖子点赞列表", func(t *testing.T) {
		// 创建多个用户点赞
		user2 := createTestUser("like_test_user2")
		user3 := createTestUser("like_test_user3")

		likeService.LikePost(post.ID, user.ID)
		likeService.LikePost(post.ID, user2.ID)
		likeService.LikePost(post.ID, user3.ID)

		likes, total, err := likeService.GetPostLikes(post.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, likes, 3)
	})

	t.Run("测试获取用户点赞的帖子列表", func(t *testing.T) {
		// 创建多个帖子
		post2 := createTestPost(user.ID, "测试帖子2", "这是第二个测试帖子")
		post3 := createTestPost(user.ID, "测试帖子3", "这是第三个测试帖子")

		// 用户点赞多个帖子
		likeService.LikePost(post.ID, user.ID)
		likeService.LikePost(post2.ID, user.ID)
		likeService.LikePost(post3.ID, user.ID)

		likes, total, err := likeService.GetUserLikedPosts(user.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, likes, 3)
	})

	t.Run("测试获取帖子点赞数", func(t *testing.T) {
		// 创建多个用户点赞
		user4 := createTestUser("like_test_user4")
		user5 := createTestUser("like_test_user5")

		likeService.LikePost(post.ID, user.ID)
		likeService.LikePost(post.ID, user4.ID)
		likeService.LikePost(post.ID, user5.ID)

		count, err := likeService.GetPostLikeCount(post.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	// 清理测试数据
	cleanupTestDB()
}

// TestLikeAPI 测试点赞API接口
func TestLikeAPI(t *testing.T) {
	// 这里可以添加HTTP API测试
	// 需要启动服务器并发送HTTP请求
}

// BenchmarkLikeService 性能测试
func BenchmarkLikeService(b *testing.B) {
	b.Skip("需要完整的应用环境，跳过性能测试")
	setupTestDB()
	defer cleanupTestDB()

	user := createTestUser("benchmark_user")
	post := createTestPost(user.ID, "基准测试帖子", "用于性能测试的帖子")
	likeService := appService.LikeServer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		likeService.LikePost(post.ID, user.ID)
		likeService.UnlikePost(post.ID, user.ID)
	}
}
