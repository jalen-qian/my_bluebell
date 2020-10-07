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

// AddPostHandler 创建帖子接口处理函数
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
	p.UserId = getCurrentUserId(c)
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

// PostListHandler 帖子列表接口处理函数
func PostListHandler(c *gin.Context) {
	// 1.获取分页参数
	page, pageSize := getPageParams(c)
	zap.L().Info("get Page and PageSize", zap.Int("page", page), zap.Int("pageSize", pageSize))
	// 2.从数据库中查询
	postList, err := logic.GetPostList(page, pageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList(page,pageSize) failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeUsualFailed, "获取帖子列表失败")
		return
	}
	ResponseSuccess(c, postList)
}

// PostList2Handler 帖子列表接口优化版
// 支持根据评分和发布时间排序
func PostList2Handler(c *gin.Context) {
	//定义参数，默认第一页，一页10条数据，按照发布时间排序
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: "time",
	}
	//参数绑定与参数校验
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostList2Handler with invalid params", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("PostList2Handler params", zap.Any("params", p))

	//从逻辑层获取数据
	postList, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList(page,pageSize) failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeUsualFailed, "获取帖子列表失败")
		return
	}
	ResponseSuccess(c, postList)

}

// PostVoteHandler 帖子投票处理函数
func PostVoteHandler(c *gin.Context) {
	// 1.参数校验 文章ID 投票情况(1赞成票 -1反对票 0取消之前的投票)
	params := new(models.ParamPostVote)
	if err := c.ShouldBindJSON(params); err != nil {
		//类型断言，判断是否是因为validator出错
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		//不是validator报错，说明是其他参数错误
		ResponseError(c, CodeInvalidParam)
		return
	}
	userId := getCurrentUserId(c)
	// 2.业务逻辑处理
	if err := logic.PostVote(userId, params); err != nil {
		zap.L().Error("logic.PostVote(userId, params) failed",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Any("params", params))
		//帖子不存在 或者 投票期已过，返回对应的错误
		if errors.Is(err, logic.ErrPostNotExists) || errors.Is(err, logic.ErrVoteTimeOut) {
			ResponseErrorWithMsg(c, CodeUsualFailed, err.Error())
			return
		}
		ResponseErrorWithMsg(c, CodeUsualFailed, "投票失败")
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}
