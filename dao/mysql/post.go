package mysql

import (
	"bluebell/models"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
                 post_id, title, content, author_id, community_id)
                 values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
			post_id, title, content, author_id, community_id, create_time
			from post 
			where post_id = ?
			`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表
func GetPostList(PageNum, PageSize int64) (posts []*models.Post, err error) {
	sqlStr := `select 
			post_id, title, content, author_id, community_id, create_time
			from post 
			order by create_time
			desc
			limit ?,?
			`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (PageNum-1)*PageSize, PageSize)
	if err != nil {
		zap.L().Error("db.Get(posts, sqlStr) failed", zap.Error(err))
	}
	return posts, err
}
