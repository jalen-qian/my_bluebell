package controller

import (
	"errors"
	"strconv"

	"github.com/bluebell/dao/mysql"

	"github.com/bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区列表接口
func CommunityHandler(c *gin.Context) {
	//调用Logic层，查询所有的社区列表，返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeUsualFailed, "获取社区列表失败")
		return
	}
	zap.L().Info("CommunityList", zap.Any("data", data))
	ResponseSuccess(c, data)
}

// 社区详情接口
func CommunityDetailHandler(c *gin.Context) {
	// 获取参数中的社区ID
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从数据库中查询社区详情
	commDetail, err := logic.GetCommunityDetail(int64(id))
	if err != nil {
		if errors.Is(err, mysql.ErrCommunityNotFound) {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseErrorWithMsg(c, CodeUsualFailed, "获取社区详情失败")
		return
	}
	ResponseSuccessWithMsg(c, commDetail, "获取社区详情成功")
}
