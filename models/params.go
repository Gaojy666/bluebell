package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求的参数结构体
// ParamSignUp注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户,这里不需要写
	PostID    string `json:"post_id" binding:"required"`              // 帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 -1 0"` // 赞成票(1)or反对票(-1)取消投票(0)
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	// CommunityID 可以为空
	CommunityID int64  `json:"community_id" from:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数据量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

//// ParamCommunityPostList 按社区获取帖子列表query string参数
//type ParamCommunityPostList struct {
//	ParamPostList
//}
