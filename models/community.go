package models

import (
	"time"

	"gorm.io/gorm"
)

// 社区表 每篇帖子都归属于一个唯一的社区，一个用户可以加入到多个社区
type Community struct {
	CommunityId   int64      `json:"id" gorm:"NOT NULL;COMMENT:社区ID; uniqueIndex:idx_community_id"`
	CommunityName string     `json:"name" gorm:"type:varchar(128); NOT NULL; DEFAULT:''; COMMENT:社区名称; uniqueIndex:idx_community_name"`
	Introduction  string     `json:"introduction" gorm:"type:varchar(256); NOT NULL; DEFAULT:''; COMMENT:详情"`
	gorm.Model    `json:"-"` //忽略掉
}

// 定义社区详情结构体，对应社区表中的3个字段，并规定json返回的字段
type CommunityDetail struct {
	Id         int64     `json:"id" gorm:"column:community_id"`
	Name       string    `json:"name" gorm:"column:community_name"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_at"`
}
