package frontend

import (
	"testing"

	"server/model/frontend"

	"github.com/stretchr/testify/assert"
)

func TestArticleServiceStruct(t *testing.T) {
	svc := &Article{}
	assert.NotNil(t, svc)
}

func TestEncodeDecodeArticleListCache(t *testing.T) {
	t.Run("编码后解码应恢复原数据", func(t *testing.T) {
		original := []frontend.Article{
			{
				Title:   "文章1",
				Content: "内容1",
				UserId:  1,
			},
			{
				Title:   "文章2",
				Content: "内容2",
				UserId:  2,
			},
		}

		data, err := encodeArticleListCache(original)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)

		decoded, err := decodeArticleListCache(data)
		assert.NoError(t, err)
		assert.Len(t, decoded, 2)
		assert.Equal(t, "文章1", decoded[0].Title)
		assert.Equal(t, "内容1", decoded[0].Content)
		assert.Equal(t, "文章2", decoded[1].Title)
		assert.Equal(t, "内容2", decoded[1].Content)
	})

	t.Run("编码空列表", func(t *testing.T) {
		original := []frontend.Article{}

		data, err := encodeArticleListCache(original)
		assert.NoError(t, err)

		decoded, err := decodeArticleListCache(data)
		assert.NoError(t, err)
		assert.Len(t, decoded, 0)
	})

	t.Run("解码无效数据应返回错误", func(t *testing.T) {
		_, err := decodeArticleListCache([]byte("invalid data"))
		assert.Error(t, err)
	})
}

func TestEncodeDecodeTotalCache(t *testing.T) {
	t.Run("编码后解码应恢复原值", func(t *testing.T) {
		original := int64(12345)
		data := encodeTotalCache(original)
		assert.Len(t, data, 8)
		assert.Equal(t, original, decodeTotalCache(data))
	})

	t.Run("短期数据应返回0", func(t *testing.T) {
		assert.Equal(t, int64(0), decodeTotalCache([]byte{1, 2, 3}))
	})

	t.Run("编码零值", func(t *testing.T) {
		data := encodeTotalCache(0)
		assert.Len(t, data, 8)
		assert.Equal(t, int64(0), decodeTotalCache(data))
	})
}

func TestGetArticleList(t *testing.T) {
	t.Skip("需要完整的应用环境（REDIS + DB），跳过")
}

func TestGetArticleDetail(t *testing.T) {
	t.Skip("需要完整的应用环境（REDIS + DB），跳过")
}

func TestGetSearchArticle(t *testing.T) {
	t.Skip("需要完整的应用环境（DB），跳过")
}
