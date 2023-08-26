package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(usersID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", usersID),
		zap.String("postID", p.PostID),
		zap.Int8("dirction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(usersID)), p.PostID, float64(p.Direction))
}
