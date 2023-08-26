package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	orderTime  = "time"
	orderScore = "score"
)

// CreatePostHandler 创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		// zap.Any("err", err) 可以用于记录各种类型的字段
		// 而 zap.Error(err) 则专门用于记录错误信息，使得日志中的错误更加突出和易于识别。
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c取到当前发请求的用户id
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数(帖子的id，从URL中获取)
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id取出帖子规模
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListDetail 获取帖子列表的接口
func GetPostListHandler(c *gin.Context) {
	// 帖子列表有可能上万条，无法全部展示，可以分页展示
	// 获取分页参数
	PageNum, PageSize := GetPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(PageNum, PageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListDetail2 根据时间或分数排序
// 根据前端传来的参数，动态的获取帖子列表
// 按创建时间排序，或者 按照分数排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询详细信息
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: orderTime, // magic string 防止硬编码
	}
	if err := c.ShouldBindQuery(); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// c.shouldBind()  根据请求的数据选择相应的方法去获取数据
	// c.ShouldBindJson()// 如果请求中携带的是Json格式的数据，采用这个方法获取到数据

	// 获取分页参数
	PageNum, PageSize := GetPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(PageNum, PageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
