package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取分页参数
func getPageParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	return page, pageSize
}

func getCurrentUserId(c *gin.Context) int64 {
	return c.GetInt64(ContextUserIdKey)
}
