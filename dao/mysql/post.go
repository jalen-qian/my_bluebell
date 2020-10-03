package mysql

import (
	"errors"

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
