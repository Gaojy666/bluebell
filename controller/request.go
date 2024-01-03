package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIdKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUserID 获取当前登录的用户ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	// 接口值转换为int64 ????????????????
	uid, ok := c.Get(CtxUserIdKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPageInfo(c *gin.Context) (int64, int64) {
	PageNumStr := c.Query("page")
	PageSizeStr := c.Query("size")
	var (
		PageNum  int64
		PageSize int64
		err      error
	)
	PageNum, err = strconv.ParseInt(PageNumStr, 10, 64)
	if err != nil {
		PageNum = 1
	}
	PageSize, err = strconv.ParseInt(PageSizeStr, 10, 64)
	if err != nil {
		PageSize = 10
	}
	return PageNum, PageSize
}
