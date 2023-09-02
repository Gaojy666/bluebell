package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "lu741208",
		DbName:       "bluebell",
		Port:         3306,
		MaxOpenConn:  200,
		MaxIdleConns: 50,
	}
	err := InitDB(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("TestCreatePost insert record mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record mysql succeeded")
}
