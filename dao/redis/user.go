package redis

import (
	"fmt"
	"time"
)

// GetUserToken 通过用户ID从Redis中取一个持久化的Token
func GetUserToken(userId int64) (token string, err error) {
	return rdb.Get(getUserTokenKey(userId)).Result()
}

// SetUserToken 将一个用户Token绑定到
func SetUserToken(userId int64, token string, expiration time.Duration) error {
	return rdb.Set(getUserTokenKey(userId), token, expiration).Err()
}

func getUserTokenKey(userId int64) string {
	return fmt.Sprintf("%s%d", KeyUserTokenPF, userId)
}
