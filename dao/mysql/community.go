package mysql

import (
	"errors"

	"go.uber.org/zap"

	"github.com/bluebell/models"
	"gorm.io/gorm"
)

func GetCommunityList() ([]models.Community, error) {
	var communities = make([]models.Community, 0, 0)
	result := Db.Find(&communities)
	if result.Error != nil {
		return communities, result.Error
	}
	return communities, nil
}

// 通过社区ID查询社区
func GetCommunityDetailById(id int64) (communityDetail *models.CommunityDetail, err error) {
	//查询
	communityDetail = new(models.CommunityDetail)
	result := Db.Model(&models.Community{}).Select("community_id", "community_name", "introduction", "created_at").
		Where("community_id = ?", id).
		First(communityDetail)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//返回数据不存在
			return nil, ErrCommunityNotFound
		}
		return nil, result.Error
	}
	zap.L().Debug("查询到社区信息", zap.Any("communityDetail", communityDetail))
	return communityDetail, nil
}

// IsCommunityExists 判断社区是否存在
func IsCommunityExists(communityId int64) (isExist bool, err error) {
	c := new(models.Community)
	result := Db.Model(c).Select("community_id").Where("community_id = ?", communityId).First(c)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
