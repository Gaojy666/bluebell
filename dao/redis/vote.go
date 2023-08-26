package redis

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

// 投票功能
// 1.用户投票的数据

// 本项目使用简单版的投票分数
// 投一票就加432分 86400s/200 -> 需要200张赞成票可以给你的帖子续一天 -> <redis实战>

/* 投票的几种情况
dirction=1时，有两种情况：
	1.之前没有投过票，现在投赞成票   差值的绝对值：1  +432
	2.之前投反对票，现在改投赞成票   差值的绝对值：2  +432*2
dirction=0时，有两种情况：
	1.之前投过赞成票，现在要取消投票  差值的绝对值：1  -432
	2.之前投过反对票，现在要取消投票  差值的绝对值：1  +432
dirction=-1时，有两种情况：
	1.之前没有投过票，现在投反对票  差值的绝对值：1  -432
	2.之前投赞成票，现在改投反对票  差值的绝对值：2  -432*2

投票的限制：
防止挖坟！！！
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了
	1.到期之后，将redis中保存的赞成票数及反对票数存储到mysql表中(持久化)
	2.到期之后，删除那个KeyPostVotedZSetPrefix
*/

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZset), postID).Val()
	// 每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2.更新帖子的分数
	// 先查之前的投票纪录
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 查询当前用户给这个帖子投票的分数
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), userID).Val()
	var flag float64
	if value > ov {
		flag = 1
	} else {
		flag = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 更新分数
	_, err := client.ZIncrBy(ctx, getRedisKey(KeyPostScoreZset), flag*diff*scorePerVote, postID).Result()
	if ErrVoteTimeExpire != nil {
		return err
	}
	// 3. 记录用户为该帖子投过票
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if value == 0 {
		// 取消投票
		_, err = client.ZRem(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), userID).Result()
	} else {
		_, err = client.ZAdd(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), &redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		}).Result()
	}
	return err
}
