package redis

import (
	"github.com/bluebell/models"
	"go.uber.org/zap"
)

// 通过参数获取帖子Id列表
func GetPostIdsByPostListParams(p *models.ParamPostList) (ids []string, err error) {
	key := KeyPostTime
	if p.Order == models.OrderTypeScore {
		key = KeyPostScore
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	zap.L().Debug("GetPostIdsByPostListParams start and end",
		zap.Int64("start", start),
		zap.Int64("end", end))

	return rdb.ZRevRange(key, start, end).Result()
}
