package logic

import "github.com/bluebell/dao/mysql"

// 所有社区相关的业务逻辑，写在这里

// GetCommunityList 获取所有的社区列表
// 返回字段包括社区Id，社区名称，详情
func GetCommunityList() (data interface{}, err error) {
	return mysql.GetCommunityList()
}
