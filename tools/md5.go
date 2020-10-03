package tools

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/bluebell/settings"
)

/**
  Md5Encrypt 对一个字符串进行md5加密
*/
func Md5Encrypt(str string) string {
	h := md5.New()
	h.Write([]byte(settings.Conf.Md5Secret))
	return hex.EncodeToString(h.Sum([]byte(str)))
}
