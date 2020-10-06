package mysql

import (
	"errors"

	"github.com/bluebell/tools"
	"gorm.io/gorm"

	"github.com/bluebell/models"
)

// CheckUserExists 检查用户名对应的用户是否存在
func CheckUserExists(userName string) (isExist bool, err error) {
	u := &models.User{UserName: userName}
	result := Db.Where("user_name = ?", userName).First(u)
	if result.Error != nil {
		//如果是因为用户未找到，业务上不能算出现错误，只能算用户不存在
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		//其他情况下，算数据库查询出错，这时候应该返回错误
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// SaveUser 保存一个新的用户
func SaveUser(user *models.User) (err error) {
	err = Db.Save(user).Error
	if err != nil {
		return err
	}
	return
}

// Login 登录
func Login(user *models.User) (err error) {
	iptPassword := tools.Md5Encrypt(user.Password)

	result := Db.Select("user_id", "user_name", "password").
		Where("user_name = ?", user.UserName).
		First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//返回用户不存在的错误
			return ErrUserNotExist
		}
		return ErrServerBusy
	}
	//判断输入的密码和数据库中的密码是否一致
	if iptPassword != user.Password {
		return ErrInvalidPassword
	}
	return
}

// 根据用户ID查询用户信息
func GetUserById(uId int64) (user *models.User, err error) {
	user = new(models.User)
	result := Db.Model(user).
		Select("user_id", "user_name").
		Where("user_id = ?", uId).
		First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotExist
		}
		return nil, result.Error
	}
	return user, nil
}
