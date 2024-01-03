package redis

import (
	"context"
	"strconv"
	"time"
)

func StoreUserToken(userID int64, token string) error {
	// 存储user_id和token对应关系
	ukey := getRedisKey(KeyUserToken) + strconv.Itoa(int(userID))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.Set(ctx, ukey, token, 0).Result()
	return err
}

func GetTokenFromID(userID int64) (string, error) {
	// 存储user_id和token对应关系
	ukey := getRedisKey(KeyUserToken) + strconv.Itoa(int(userID))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return client.Get(ctx, ukey).Result()
}
