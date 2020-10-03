package mysql

import "github.com/bluebell/models"

func GetCommunityList() ([]models.Community, error) {
	var communities = make([]models.Community, 0, 0)
	result := Db.Find(&communities)
	if result.Error != nil {
		return communities, result.Error
	}
	return communities, nil
}
