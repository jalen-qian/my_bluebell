package redis

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/bluebell/models"
	"github.com/go-redis/redis"
)

// 获取帖子的发布时间
func GetPostPublishTime(postId int64) (pubTime int64, err error) {
	postIdStr := strconv.FormatInt(postId, 10)
	ret := rdb.ZScore(KeyPostTime, postIdStr)
	err = ret.Err()
	if err != nil {
		return 0, err
	}
	return int64(ret.Val()), err
}

// 保存帖子到Redis（帖子发布时间 & 初始化帖子分数）
func SavePost(postId int64) (err error) {
	pipe := rdb.TxPipeline()
	//保存帖子的发布时间
	pipe.ZAdd(KeyPostTime, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	pipe.ZAdd(KeyPostScore, redis.Z{
		Score:  0,
		Member: postId,
	})
	_, err = pipe.Exec()
	return
}

// 给文章投票
func VoteForPost(userId int64, p *models.ParamPostVote) (err error) {
	//查询之前的投票情况
	keyPostVoted := fmt.Sprintf("%s%d", KeyPostVotedPF, p.PostId)
	oldVote := rdb.ZScore(keyPostVoted, strconv.FormatInt(userId, 10)).Val()
	pipe := rdb.TxPipeline()

	//投票分数差 = (当前的投票 - 之前的投票) * 432
	direction, _ := p.GetFloat64Direction()
	scoreDiff := (direction - oldVote) * 432
	zap.L().Debug("scoreDiff",
		zap.Float64("scoreDiff", scoreDiff),
		zap.Float64("oldVote", oldVote))
	//如果投票分数不为0，则投票
	if scoreDiff != 0 {
		pipe.ZIncrBy(KeyPostScore, scoreDiff, strconv.FormatInt(p.PostId, 10))
	}
	if direction == 1 || direction == -1 {
		pipe.ZAdd(keyPostVoted, redis.Z{
			Score:  direction,
			Member: userId,
		})
	} else {
		pipe.ZRem(keyPostVoted, userId)
	}
	_, err = pipe.Exec()
	return
}
