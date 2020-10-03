package logic

import (
	"github.com/bluebell/dao/mysql"
	"github.com/bluebell/dao/redis"
	"github.com/bluebell/models"
	"github.com/bluebell/pkg/jwt"
	"github.com/bluebell/pkg/snowflake"
	"github.com/bluebell/tools"
	"go.uber.org/zap"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断这个用户是否存在
	user := models.User{
		UserName: p.Username,
	}
	var isUserExists bool
	if isUserExists, err = mysql.CheckUserExists(p.Username); err != nil {
		//查询出错，直接返回错误
		return err
	}
	//查询没有出错，但用户已存在
	if isUserExists {
		return mysql.ErrUserExist
	}
	//用户不存在, 保存用户进入数据库
	user.Password = tools.Md5Encrypt(p.Password)
	user.UserId, err = snowflake.GetId()
	if err != nil {
		//雪花算法出错，也直接返回错误
		return err
	}
	err = mysql.Db.Save(&user).Error
	if err != nil {
		return err
	}

	return
}

// 登录逻辑
func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		UserName: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return "", err
	}
	zap.L().Debug("登录成功，用户信息：",
		zap.Int64("userId", user.UserId),
		zap.String("userName", user.UserName))
	//登录成功，这里由于传的是指针，已经拿到了userId和userName
	token, err = jwt.GenToken(user.UserId, user.UserName)
	if err != nil {
		return "", err
	}
	//将Token存储到Redis中，实现单点登录，同一时刻只能有一个Token是有效的
	err = redis.SetUserToken(user.UserId, token, jwt.TokenExpireDuration)
	// 测试，将过期设置为20秒，观察20秒后能否退出登录
	//err = redis.SetUserToken(user.UserId, token, time.Second*20)
	if err != nil {
		return "", err
	}
	return token, err
}
