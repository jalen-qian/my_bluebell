package logic

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bluebell/dao/redis"

	"github.com/bluebell/dao/mysql"
	"github.com/bluebell/models"
	"github.com/bluebell/pkg/snowflake"
)

const (
	VoteTimeLimit = 604800 //投票截止时间限制（一周的秒数）
)

var (
	ErrVoteTimeOut   = errors.New("投票时间已过")
	ErrPostNotExists = errors.New("帖子不存在")
)

// SavePost 保存帖子业务逻辑处理
func SavePost(p *models.ParamAddPost) (err error) {
	//判断社区是否存在
	isExist, err := mysql.IsCommunityExists(p.CommunityId)
	fmt.Printf("帖子是否存在：%v,%v\n", isExist, err)
	if err != nil {
		return err
	}
	//社区不存在，报错
	if !isExist {
		return mysql.ErrCommunityNotFound
	}
	postId, _ := snowflake.GetId()
	//保存帖子到数据库
	post := &models.Post{
		PostId:      postId,
		AuthorId:    p.UserId,
		CommunityId: p.CommunityId,
		Title:       p.Title,
		Content:     p.Content,
	}
	if err = mysql.SavePost(post); err != nil {
		return err
	}
	//保存成功后，在Reids记录发布时间
	err = redis.SavePost(postId)
	return
}

// 业务逻辑处理：根据文章ID查询文章详细信息
func GetPostDetailByPostId(postId int64) (data *models.ApiPostDetail, err error) {
	// 1.根据文章Id查询文章信息
	post, err := mysql.GetPostDetailByPostId(postId)
	if err != nil {
		return nil, err
	}
	// 2. 根据文章中的CommunityId查询Community相关信息
	commDetail, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		return nil, err
	}
	// 3.根据文章中的author_id查询用户相关信息
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		return nil, err
	}
	data = &models.ApiPostDetail{
		Post:            post,
		AuthorName:      user.UserName,
		CommunityDetail: commDetail,
	}
	return
}

// GetPostList 分页查询帖子列表的业务逻辑
func GetPostList(page, pageSize int) (data []*models.ApiPostDetail, err error) {
	//1.分页查询帖子列表
	posts, err := mysql.GetPostList(page, pageSize)
	if err != nil {
		return nil, err
	}
	//2.遍历每个帖子，查询帖子包含的社区数据和用户名称
	data = fillCommunityAndAuthorForPosts(&posts)
	return
}

// 遍历每个帖子，查询帖子包含的社区数据和用户名称
func fillCommunityAndAuthorForPosts(posts *[]*models.Post) []*models.ApiPostDetail {
	data := make([]*models.ApiPostDetail, 0, len(*posts))
	for _, post := range *posts {
		commDetail, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			continue
		}
		// 3.根据文章中的author_id查询用户相关信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			continue
		}
		data = append(data, &models.ApiPostDetail{
			Post:            post,
			AuthorName:      user.UserName,
			CommunityDetail: commDetail,
		})
	}
	return data
}

// 获取帖子列表接口2
// 可以根据评分或者发布时间来排序
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//首先从redis读取文章ID
	ids, err := redis.GetPostIdsByPostListParams(p)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("redis.GetPostIdsByPostListParams(p)", zap.Any("ids", ids), zap.Any("params", p))
	//将这些ID转换为int64类型
	//通过Mysql查询这些ID对应的帖子
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIds(ids) failed", zap.Error(err))
		return nil, err
	}
	//对每个帖子按照ids的顺序排序

	//2.遍历每个帖子，查询帖子包含的
	data = fillCommunityAndAuthorForPosts(&posts)
	return data, nil
}

/**
PostVote 帖子评分
投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

投票分3种情况 1-赞成票 0取消投票 -1反对票
	情况1：投赞成票(1)				投票差值   分数差值
		a.之前没有投过票，现在投赞成票  1          432
		b.之前投过反对票，现在投赞成票  2          432 * 2

	情况2：取消投票(0)
		a.之前投赞成票，现在取消投票	-1		    -432
		b.之前投反对票，现在取消投票	1			432
		c.之前没有投过票，什么都不用做

	情况3：投反对票
		a.之前没有投过票，现在投反对票  -1          -432
		b.之前投过赞成票，现在投反对票  -2          -432 * 2

   投票的限制：
   每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
   	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
   	2. 到期之后删除那个 keyPostVotedPF
*/
func PostVote(userId int64, p *models.ParamPostVote) (err error) {
	//首先拿到帖子发帖日期，看当前时间是否已过
	pubTime, err := redis.GetPostPublishTime(p.PostId)
	if err != nil {
		zap.L().Error("redis.GetPostPublishTime(p.PostId) got err",
			zap.Error(err),
			zap.Int64("pubTime", pubTime),
			zap.Int64("userId", userId),
			zap.Int64("postId", p.PostId))
		return ErrPostNotExists
	}
	if time.Now().Unix()-pubTime > VoteTimeLimit {
		return ErrVoteTimeOut
	}
	//然后投票
	err = redis.VoteForPost(userId, p)
	return
}
