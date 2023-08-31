package mysql

import (
	"bluebell/models"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
                 post_id, title, content, author_id, community_id)
                 values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// 根据ID查询单个帖子详情
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

// GetPostListByIDs 根据给定的ID列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	// FIND_IN_SET，根据指定顺序返回，否则默认按照id顺序返回
	sqlStr := `select
			post_id, title, content, author_id, community_id, create_time
			from post
			while post_id  in (?)
			order by FIND_IN_SET(post_id, ?)
		`
	// sqlx.In批量查询
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	fmt.Printf("query: %#v\n", query)
	fmt.Printf("args: %#v\n", args)
	if err != nil {
		zap.L().Error("db.Get(posts, sqlStr) failed", zap.Error(err))
		return nil, err
	}

	query = db.Rebind(query)
	fmt.Printf("query: %#v\n", query)

	err = db.Select(&postList, query, args...)
	return
}
