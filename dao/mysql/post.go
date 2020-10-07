package mysql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bluebell/models"
	"gorm.io/gorm"
)

// 数据库层保存帖子处理
func SavePost(post *models.Post) (err error) {
	err = Db.Save(post).Error
	return
}

// 根据帖子ID查询帖子
func GetPostDetailByPostId(postId int64) (post *models.Post, err error) {
	post = new(models.Post)
	result := Db.Model(post).Select("post_id", "title", "author_id", "community_id", "content", "status", "created_at", "updated_at").Where("post_id = ?", postId).First(post)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCommunityNotFound
		}
		return nil, result.Error
	}
	return
}

//分页查询文章列表
func GetPostList(page, pageSize int) (postList []*models.Post, err error) {
	//先分配内存
	postList = make([]*models.Post, 0, pageSize)
	//查询
	err = Db.Model(&models.Post{}).
		Select("post_id", "title", "author_id", "community_id", "content", "status", "created_at", "updated_at").
		Where("id > ? ", 0).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&postList).Error
	if err != nil {
		//出错，直接返回
		return
	}
	//没有出错，说明查询出来了，直接返回
	return
}

// 通过帖子ID列表查询帖子
func GetPostListByIds(postIds []string) (postList []*models.Post, err error) {
	//先分配内存
	postList = make([]*models.Post, 0, len(postIds))
	//查询
	ids := make([]int64, len(postIds))
	for i, id := range postIds {
		ids[i], _ = strconv.ParseInt(id, 0, 64)
	}
	err = Db.Model(&models.Post{}).
		Select("post_id", "title", "author_id", "community_id", "content", "status", "created_at", "updated_at").
		Where("post_id in(?)", ids).
		Order(fmt.Sprintf("FIND_IN_SET(post_id,'%s')", strings.Join(postIds, ","))).
		Find(&postList).Error
	if err != nil {
		//出错，直接返回
		return
	}
	//没有出错，说明查询出来了，直接返回
	return
}
