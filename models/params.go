package models

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
	PostID    string `json:"post_id" binding:"required"`                       // 帖子id
	Direction int8   `json:"direction,string" binding:"required,oneof=1 -1 0"` // 赞成票(1)or反对票(-1)取消投票(0)
}
