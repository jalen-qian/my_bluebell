package controller

type ResCode int64

//定义一些常用的固定状态码
const (
	CodeSuccess         ResCode = 1000 + iota //请求成功
	CodeUsualFailed                           //通用失败场景（某些未知名错误等）
	CodeInvalidParam                          //请求参数有误
	CodeUserExist                             //用户已存在
	CodeUserNotExist                          //用户不存在
	CodeInvalidPassword                       //密码错误
	CodeServerBusy                            //服务器繁忙

	CodeNeedLogin    //需要登录
	CodeInvalidToken //无效的token

	CodeCommunityNotExist //社区不存在
	CodePostNotExist      //帖子不存在
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:           "success",
	CodeUsualFailed:       "请求失败", //通用场景直接返回请求失败
	CodeInvalidParam:      "请求参数错误",
	CodeUserExist:         "用户已存在",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "密码错误",
	CodeServerBusy:        "服务繁忙",
	CodeNeedLogin:         "需要登录",
	CodeInvalidToken:      "无效的Token",
	CodeCommunityNotExist: "社区不存在",
	CodePostNotExist:      "帖子不存在",
}

func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeServerBusy.Msg()
	}
	return msg
}
