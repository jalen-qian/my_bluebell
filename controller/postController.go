package controller

import (
	"errors"
	"strconv"

	"github.com/bluebell/dao/mysql"
	"github.com/bluebell/logic"
	"github.com/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

/**
在这个文件中定义帖子相关接口的处理函数
*/

// 创建帖子
func AddPostHandler(c *gin.Context) {
	// 1.获取参数 标题 内容 社区ID并校验参数
	p := new(models.ParamAddPost)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("创建帖子 c.ShouldBindJSON(p) failed", zap.Error(err))
		//记录参数校验错误的日志
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, CodeInvalidParam)
	}
	// 2.保存帖子
	p.UserId = generateCurrentUserId(c)
	if err := logic.SavePost(p); err != nil {
		zap.L().Error("创建帖子 logic.SavePost(p) failed", zap.Error(err))
		if errors.Is(err, mysql.ErrCommunityNotFound) {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeUsualFailed)
		return
	}
	// 3.返回响应
	ResponseSuccessWithMsg(c, nil, "创建帖子成功")
}

// PostDetailHandler 帖子详情接口处理函数
func PostDetailHandler(c *gin.Context) {
	//1.参数校验，获取帖子ID
	strPostId := c.Param("postId")
	postId, err := strconv.Atoi(strPostId)
	if err != nil {
		zap.L().Error("PostDetail Invalid Params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据获取帖子详情
	post, err := logic.GetPostDetailByPostId(int64(postId))
	if err != nil {
		zap.L().Error("logic.GetPostDetail failed", zap.Error(err))
		if errors.Is(err, mysql.ErrPostNotFound) {
			ResponseError(c, CodePostNotExist)
			return
		}
		ResponseError(c, CodeUsualFailed)
		return
	}
	//3.返回正确响应
	ResponseSuccess(c, post)
}
