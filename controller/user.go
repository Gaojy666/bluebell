package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	//fmt.Println("0-------------")

	if err := c.ShouldBindJSON(p); err != nil {
		// 检查错误是否是validator的错误，如果是，执行翻译
		// 否则非validator.ValidationErrors类型错误直接返回, 比如传过来的字段都对应不上
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动对求情参数进行详细的业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	//fmt.Println(p)
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		//zap.L().Error("--------", zap.Error(err))
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 检查错误是否是validator的错误，如果是，执行翻译
		// 否则非validator.ValidationErrors类型错误直接返回, 比如传过来的字段都对应不上
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)), // 翻译错误
		})
		return
	}
	// 2.业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // 如果id值1<<53-1，int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}
