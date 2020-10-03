package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   int64  `gorm:"NOT NULL;COMMENT:用户ID; uniqueIndex:idx_user_id"`
	UserName string `gorm:"type:varchar(64); NOT NULL; DEFAULT:''; COMMENT:用户名; uniqueIndex:idx_username"`
	Password string `gorm:"type:varchar(64); NOT NULL; DEFAULT:''; COMMENT:密码"`
	Email    string `gorm:"type:varchar(64); NOT NULL; DEFAULT:''; COMMENT:邮箱"`
	Gender   int    `gorm:"type:tinyint(4); NOT NULL; DEFAULT:0; COMMENT:性别"`
}
