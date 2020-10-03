package logic

import (
	"fmt"

	"github.com/bluebell/dao/mysql"
	"github.com/bluebell/models"
	"github.com/bluebell/pkg/snowflake"
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
