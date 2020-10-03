package models

import (
	"time"

	"gorm.io/gorm"
)

// 帖子（文章）Model 注意内存对其，可以使结构体对象的内存占用更小
type Post struct {
	ID          uint           `json:"-"gorm:"primarykey;COMMENT:主键ID"`
	PostId      int64          `json:"post_id"gorm:"type:bigint(20);NOT NULL;COMMENT:帖子ID;uniqueIndex:idx_post_id"`
	AuthorId    int64          `json:"author_id"gorm:"type:bigint(20);NOT NULL;DEFAULT:0;COMMENT:作者的用户id;index:idx_author_id"`
	CommunityId int64          `json:"community_id"gorm:"type:bigint(20);NOT NULL;DEFAULT:0;COMMENT:所属社区;index:idx_community_id"`
	Status      uint8          `json:"status"gorm:"type:tinyint(4);NOT NULL;DEFAULT:1;COMMENT:状态"`
	Title       string         `json:"title"gorm:"type:varchar(128);NOT NULL;DEFAULT:'';COMMENT:标题;"`
	Content     string         `json:"content"gorm:"type:varchar(8192);NOT NULL;DEFAULT:'';COMMENT:内容"`
	CreatedAt   time.Time      `json:"created_at"gorm:"COMMENT:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at"gorm:"COMMENT:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-"gorm:"index;COMMENT:删除时间"`
}

type ApiPostDetail struct {
	*Post
	AuthorName       string `json:"author_name"`
	*CommunityDetail `json:"community"`
}
