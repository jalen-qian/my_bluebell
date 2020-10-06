package jwt

import (
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/bluebell/settings"

	"github.com/dgrijalva/jwt-go"
)

// 定义Token过期时间（2个小时）
const TokenExpireDuration = time.Hour * 24 * 30

// 定义Secret
var TokenSecret = []byte(settings.Conf.TokenSecret)

// 定义自己的JWT认证结构体，存储用户信息
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userId int64, userName string) (string, error) {
	// 定义MyClaims
	c := &MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "bluebell-jalen",
		},
	}
	//生成token，指定算法类型为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(TokenSecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 使用jwt.ParseWithClaims方法解析
	myClaims := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, myClaims, func(token *jwt.Token) (i interface{}, e error) {
		return TokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		zap.L().Debug("解析Token成功",
			zap.Int64("userId", myClaims.UserId),
			zap.String("userName", myClaims.UserName))
		return myClaims, nil
	}
	return nil, errors.New("无效的Token")

}
