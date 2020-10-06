package controller

import (
	"errors"

	"github.com/bluebell/dao/mysql"

	"github.com/bluebell/logic"

	"github.com/go-playground/validator/v10"

	"github.com/bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//定义Key名称，避免出现硬编码字符串
const (
	ContextUserIdKey   = "userId"
	ContextUserNameKey = "userName"
)

func SignUpHandler(c *gin.Context) {
	//注册流程
	//1.获取参数
	paramSignUp := new(models.ParamSignUp)
	err := c.ShouldBindJSON(paramSignUp)
	//2.参数校验
	if err != nil {
		zap.L().Error("Sign up with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			//参数错误，且是校验器返回的错误
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		//参数错误，且校验器没有校验出来（可能是参数名就有错误）
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("注册参数：", zap.Any("param", paramSignUp))
	//3.交给logic注册得到结果
	if err := logic.SignUp(paramSignUp); err != nil {
		zap.L().Error("Sign up failed ", zap.Error(err))
		if errors.Is(err, mysql.ErrUserExist) {
			//如果是用户已存在，返回用户已存在的错误
			ResponseError(c, CodeUserExist)
			return
		}
		//否则返回常规错误
		ResponseErrorWithMsg(c, CodeUsualFailed, "注册失败")
		return
	}
	ResponseSuccess(c, nil)
}

/**
 * 登录方法
 */
func LoginHandler(c *gin.Context) {
	//1.获取参数&参数校验
	paramLogin := new(models.ParamLogin)
	if err := c.ShouldBindJSON(paramLogin); err != nil {
		//记录参数校验错误的日志
		if errs, ok := err.(validator.ValidationErrors); ok {
			ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
			return
		}
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.逻辑层登录
	token, err := logic.Login(paramLogin)
	if err != nil {
		zap.L().Error("logic.Login 登录失败", zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		if errors.Is(err, mysql.ErrInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		//如果Token生成失败，也返回下面的响应
		ResponseErrorWithMsg(c, CodeUsualFailed, "登录失败")
		return
	}
	//3.返回登录结果
	ResponseSuccessWithMsg(c, gin.H{"token": token}, "登录成功")
}
