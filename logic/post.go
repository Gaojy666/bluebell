package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成Post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	return mysql.CreatePost(p)
	// 3.返回
}

// GetPostByID 根据id查询帖子详情数据
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询数据并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// 根据作者id查看作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
	}
	// 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(PageNum, PageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(PageNum, PageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		data_tmp := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, data_tmp)
	}
	return
}
