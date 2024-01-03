package redis

import (
	"bluebell/settings"
	"fmt"
	"testing"
)

// 还需要初始化db
// 初始化连接
func init() {
	dbCfg := settings.RedisConfig{
		Host:     "127.0.0.1",
		Password: "",
		Port:     6379,
		DB:       0,
		PoolSize: 100,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestStoreUserToken(t *testing.T) {
	var userId int64
	userId = 2
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0MTAwNTY4MjY3NTk5NDIxN" +
		"DQsInVzZXJuYW1lIjoibHV4aWF5dWFpIiwiZXhwIjoxNzAzNjcxODQyLCJpc3MiOiJibHVlYmVsbCJ" +
		"9.ANv1aaTbKA5zwS8H5vTSgFuLeEedMx5afa-6hO0sNJQ"
	err := StoreUserToken(userId, token)
	if err != nil {
		t.Fatalf("StoreUserToken insert record Redis failed, err:%v\n", err)
	}
	t.Logf("StoreUserToken insert record Redis succeeded")
}

func TestGetTokenFromID(t *testing.T) {
	var userId int64
	userId = 1
	token, err := GetTokenFromID(userId)
	fmt.Printf("%s", token)
	if err != nil {
		t.Fatalf("GetTokenFromID get record Redis failed, err:%v\n", err)
	}
	t.Logf("GetTokenFromID get record Redis succeeded")
}
