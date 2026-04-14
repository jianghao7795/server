package tests

import (
	"fmt"
	"testing"

	global "server/model"
	"server/model/app"
	appService "server/service/app"

	"github.com/stretchr/testify/assert"
)

// TestCommentsAndLikesIntegration 测试评论和点赞功能集成
func TestCommentsAndLikesIntegration(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过集成测试")
	// 初始化测试数据库
	setupTestDB()
	defer cleanupTestDB()

	// 创建测试用户
	user1 := createTestUser("user1")
	user2 := createTestUser("user2")
	user3 := createTestUser("user3")

	// 创建测试帖子
	post := createTestPost(user1.ID, "测试帖子", "这是用于测试评论和点赞功能的帖子内容")

	// 获取服务实例
	likeService := appService.LikeServer
	commentService := appService.CommentServer

	t.Run("测试完整的评论和点赞流程", func(t *testing.T) {
		// 1. 用户1点赞帖子
		err := likeService.LikePost(post.ID, user1.ID)
		assert.NoError(t, err)

		// 验证点赞记录
		liked, err := likeService.CheckUserLiked(post.ID, user1.ID)
		assert.NoError(t, err)
		assert.True(t, liked)

		// 验证帖子点赞数
		count, err := likeService.GetPostLikeCount(post.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 2. 用户2点赞帖子
		err = likeService.LikePost(post.ID, user2.ID)
		assert.NoError(t, err)

		// 验证点赞数增加
		count, err = likeService.GetPostLikeCount(post.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 3. 用户1创建评论
		comment1 := app.Comment{
			PostId:   post.ID,
			Content:  "这是一条测试评论",
			ParentId: 0,
			UserId:   user1.ID,
			ToUserId: 0,
		}
		err = commentService.CreateComment(&comment1)
		assert.NoError(t, err)

		// 4. 用户2回复用户1的评论（楼中楼）
		comment2 := app.Comment{
			PostId:   post.ID,
			Content:  "这是对评论的回复",
			ParentId: comment1.ID,
			UserId:   user2.ID,
			ToUserId: user1.ID,
		}
		err = commentService.CreateComment(&comment2)
		assert.NoError(t, err)

		// 5. 用户3也回复用户1的评论
		comment3 := app.Comment{
			PostId:   post.ID,
			Content:  "这是另一个回复",
			ParentId: comment1.ID,
			UserId:   user3.ID,
			ToUserId: user1.ID,
		}
		err = commentService.CreateComment(&comment3)
		assert.NoError(t, err)

		// 6. 用户3点赞帖子
		err = likeService.LikePost(post.ID, user3.ID)
		assert.NoError(t, err)

		// 验证最终点赞数
		count, err = likeService.GetPostLikeCount(post.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 7. 验证评论结构
		comments, total, err := commentService.GetCommentList(post.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total) // 3条评论
		assert.Len(t, comments, 3)

		// 8. 验证楼中楼结构
		var parentComments []app.Comment
		err = global.DB.Where("post_id = ? AND parent_id = 0", post.ID).Find(&parentComments).Error
		assert.NoError(t, err)
		assert.Len(t, parentComments, 1) // 1条父评论

		var childComments []app.Comment
		err = global.DB.Where("post_id = ? AND parent_id = ?", post.ID, comment1.ID).Find(&childComments).Error
		assert.NoError(t, err)
		assert.Len(t, childComments, 2) // 2条子评论

		// 9. 验证用户点赞的帖子列表
		likes, total, err := likeService.GetUserLikedPosts(user1.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, likes, 1)
		assert.Equal(t, post.ID, likes[0].PostId)

		// 10. 用户1取消点赞
		err = likeService.UnlikePost(post.ID, user1.ID)
		assert.NoError(t, err)

		// 验证取消点赞后的状态
		liked, err = likeService.CheckUserLiked(post.ID, user1.ID)
		assert.NoError(t, err)
		assert.False(t, liked)

		count, err = likeService.GetPostLikeCount(post.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count) // 减少到2个点赞

		// 11. 验证重复点赞
		err = likeService.LikePost(post.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经点赞过这篇帖子了")

		// 12. 验证取消不存在的点赞
		err = likeService.UnlikePost(post.ID, user1.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "还没有点赞过这篇帖子")
	})

	t.Run("测试评论的层级结构", func(t *testing.T) {
		// 创建多层嵌套评论
		post2 := createTestPost(user1.ID, "测试帖子2", "用于测试多层评论")

		// 创建父评论
		parentComment := app.Comment{
			PostId:   post2.ID,
			Content:  "父评论",
			ParentId: 0,
			UserId:   user1.ID,
			ToUserId: 0,
		}
		err := commentService.CreateComment(&parentComment)
		assert.NoError(t, err)

		// 创建子评论
		childComment := app.Comment{
			PostId:   post2.ID,
			Content:  "子评论",
			ParentId: parentComment.ID,
			UserId:   user2.ID,
			ToUserId: user1.ID,
		}
		err = commentService.CreateComment(&childComment)
		assert.NoError(t, err)

		// 创建孙评论
		grandchildComment := app.Comment{
			PostId:   post2.ID,
			Content:  "孙评论",
			ParentId: childComment.ID,
			UserId:   user3.ID,
			ToUserId: user2.ID,
		}
		err = commentService.CreateComment(&grandchildComment)
		assert.NoError(t, err)

		// 验证评论层级
		var comments []app.Comment
		err = global.DB.Where("post_id = ?", post2.ID).Find(&comments).Error
		assert.NoError(t, err)
		assert.Len(t, comments, 3)

		// 验证父评论
		var parentComments []app.Comment
		err = global.DB.Where("post_id = ? AND parent_id = 0", post2.ID).Find(&parentComments).Error
		assert.NoError(t, err)
		assert.Len(t, parentComments, 1)

		// 验证子评论
		var childComments []app.Comment
		err = global.DB.Where("post_id = ? AND parent_id = ?", post2.ID, parentComment.ID).Find(&childComments).Error
		assert.NoError(t, err)
		assert.Len(t, childComments, 1)

		// 验证孙评论
		var grandchildComments []app.Comment
		err = global.DB.Where("post_id = ? AND parent_id = ?", post2.ID, childComment.ID).Find(&grandchildComments).Error
		assert.NoError(t, err)
		assert.Len(t, grandchildComments, 1)
	})

	t.Run("测试并发点赞", func(t *testing.T) {
		post3 := createTestPost(user1.ID, "并发测试帖子", "用于测试并发点赞")

		// 模拟多个用户同时点赞
		done := make(chan bool, 10)
		for i := range 10 {
			go func(userId uint) {
				defer func() { done <- true }()
				likeService.LikePost(post3.ID, userId)
			}(uint(i + 100)) // 使用不同的用户ID
		}

		// 等待所有goroutine完成
		for range 10 {
			<-done
		}

		// 验证点赞数
		count, err := likeService.GetPostLikeCount(post3.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})
}

// TestCommentService 测试评论服务
func TestCommentService(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过服务测试")
	setupTestDB()
	defer cleanupTestDB()

	user := createTestUser("comment_user")
	post := createTestPost(user.ID, "评论测试帖子", "用于测试评论功能")
	commentService := appService.CommentServer

	t.Run("测试创建评论", func(t *testing.T) {
		comment := app.Comment{
			PostId:   post.ID,
			Content:  "测试评论内容",
			ParentId: 0,
			UserId:   user.ID,
			ToUserId: 0,
		}

		err := commentService.CreateComment(&comment)
		assert.NoError(t, err)
		assert.NotZero(t, comment.ID)

		// 验证评论已保存
		var savedComment app.Comment
		err = global.DB.Where("id = ?", comment.ID).First(&savedComment).Error
		assert.NoError(t, err)
		assert.Equal(t, post.ID, savedComment.PostId)
		assert.Equal(t, user.ID, savedComment.UserId)
		assert.Equal(t, "测试评论内容", savedComment.Content)
	})

	t.Run("测试获取评论列表", func(t *testing.T) {
		// 创建多条评论
		for i := range 5 {
			comment := app.Comment{
				PostId:   post.ID,
				Content:  fmt.Sprintf("评论内容 %d", i+1),
				ParentId: 0,
				UserId:   user.ID,
				ToUserId: 0,
			}
			commentService.CreateComment(&comment)
		}

		comments, total, err := commentService.GetCommentList(post.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total) // 包括之前的1条评论
		assert.Len(t, comments, 6)
	})

	t.Run("测试删除评论", func(t *testing.T) {
		comment := app.Comment{
			PostId:   post.ID,
			Content:  "待删除的评论",
			ParentId: 0,
			UserId:   user.ID,
			ToUserId: 0,
		}
		commentService.CreateComment(&comment)

		err := commentService.DeleteComment(comment.ID)
		assert.NoError(t, err)

		// 验证评论已删除
		var deletedComment app.Comment
		err = global.DB.Where("id = ?", comment.ID).First(&deletedComment).Error
		assert.Error(t, err) // 应该找不到记录
	})
}

// TestLikeService 测试点赞服务
func TestLikeService(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过服务测试")
	setupTestDB()
	defer cleanupTestDB()

	user := createTestUser("like_user")
	post := createTestPost(user.ID, "点赞测试帖子", "用于测试点赞功能")
	likeService := appService.LikeServer

	t.Run("测试点赞和取消点赞", func(t *testing.T) {
		// 点赞
		err := likeService.LikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 验证点赞状态
		liked, err := likeService.CheckUserLiked(post.ID, user.ID)
		assert.NoError(t, err)
		assert.True(t, liked)

		// 取消点赞
		err = likeService.UnlikePost(post.ID, user.ID)
		assert.NoError(t, err)

		// 验证取消点赞状态
		liked, err = likeService.CheckUserLiked(post.ID, user.ID)
		assert.NoError(t, err)
		assert.False(t, liked)
	})

	t.Run("测试获取点赞列表", func(t *testing.T) {
		// 创建多个用户点赞
		users := []app.User{
			createTestUser("user1"),
			createTestUser("user2"),
			createTestUser("user3"),
		}

		for _, u := range users {
			likeService.LikePost(post.ID, u.ID)
		}

		likes, total, err := likeService.GetPostLikes(post.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, likes, 3)
	})
}

// 辅助函数
func createTestUser(name string) app.User {
	user := app.User{
		Name:     name,
		Password: "test_password",
	}
	if global.DB != nil {
		global.DB.Create(&user)
	}
	return user
}

func createTestPost(userId uint, title, content string) app.Post {
	post := app.Post{
		Title:  title,
		Text:   content,
		UserId: userId,
	}
	if global.DB != nil {
		global.DB.Create(&post)
	}
	return post
}

func setupTestDB() {
	// 这里应该初始化测试数据库
	// 实际实现中需要根据项目配置来设置
	// 暂时跳过数据库初始化，因为需要完整的应用环境
}

func cleanupTestDB() {
	// 清理测试数据
	// 暂时跳过数据库清理，因为需要完整的应用环境
	if global.DB != nil {
		global.DB.Exec("DELETE FROM likes")
		global.DB.Exec("DELETE FROM comments")
		global.DB.Exec("DELETE FROM posts")
		global.DB.Exec("DELETE FROM users")
	}
}
