package controller

import (
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
