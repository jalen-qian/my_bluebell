package middlewares

import (
	"strings"

	"github.com/bluebell/settings"

	"github.com/bluebell/dao/redis"

	"github.com/bluebell/pkg/jwt"

	"github.com/bluebell/controller"
	"github.com/gin-gonic/gin"
)

// 添加用于Token校验的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有3种方式：1.放在请求头 2.放在URI 3.放在请求体
		// 我们这里规定bluebell的Token放在请求头中，并且以Bearer开头
		// 具体需要根据业务来自己决定

		// 1.拿到请求头中的验证信息
		auth := c.Request.Header.Get("Authorization")
		// 2.判断auth是否为空
		if auth == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort() //阻止后续处理函数的执行
			return
		}
		// 3.判断auth的格式是否正确
		// 正常Token格式为 <Bearer xxxx.xxxx.xxxxxx> 这里通过空格分开，取前两个
		authSplit := strings.SplitN(auth, " ", 2)
		if !(authSplit[0] == "Bearer" && len(authSplit) == 2) {
			//token格式不正确，报错：无效Token
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 4.token格式正确，校验token
		myClaims, err := jwt.ParseToken(authSplit[1])
		if err != nil {
			//如果解析失败，说明Token无效
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		//如果限制单点登录，则通过Redis
		if settings.Conf.SingleSignOn {
			//校验成功，将从Redis中取出Token进行比较
			redisToken, err := redis.GetUserToken(myClaims.UserId)
			if err != nil {
				//如果获取Redis存储的Token失败，说明Redis Token已失效，需要重新登录
				controller.ResponseError(c, controller.CodeNeedLogin)
				c.Abort()
				return
			}

			//如果Redis中的Token和传过来的Token不一致，也需要重新登录
			if authSplit[1] != redisToken {
				controller.ResponseError(c, controller.CodeNeedLogin)
				c.Abort()
				return
			}
		}
		//将Token中携带的参数，传到Context中
		c.Set(controller.ContextUserIdKey, myClaims.UserId)
		c.Set(controller.ContextUserNameKey, myClaims.UserName)
		c.Next()
	}
}
