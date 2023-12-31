package redis

// redis key

const (
	// 用:来分割命名空间,redis key注意使用命名空间的方式
	// 方便业务查询和拆分

	// KeyPrefix为通用前缀，可以快速地找到以项目名为开头的key
	Prefix                 = "bluebell:"
	KeyPostTimeZset        = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZset       = "post:score"  // zset;帖子及投票的分数
	KeyPostVotedZsetPrefix = "post:voted:" // zset;记录用户及投票类型;参数是post id
	keysCommunitySetPF     = "community:"  // set;保存每个分区下帖子的id
	KeyUserToken           = "token:"
)

func getRedisKey(key string) string {
	return Prefix + key
}
