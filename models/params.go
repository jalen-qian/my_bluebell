package models

import "strconv"

// 用户注册相关参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// 用户登录相关参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 发表帖子相关参数
type ParamAddPost struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommunityId int64  `json:"community_id" binding:"required"`
	UserId      int64  `json:"-"`
}

// 给帖子投票相关参数
type ParamPostVote struct {
	//帖子ID
	PostId int64 `json:"post_id,string" binding:"required"`
	//投票类型 1:赞成 0:取消 -1:反对
	Direction string `json:"direction" binding:"required,oneof=-1 0 1"`
}

func (p *ParamPostVote) GetFloat64Direction() (float64, error) {
	return strconv.ParseFloat(p.Direction, 0)
}
