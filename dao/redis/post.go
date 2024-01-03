package redis

import (
	"bluebell/models"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 2.确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVRANGE查询（降序）,按分数从大到小的顺序查询指定数量的元素
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 从一个有序的集合里面按照降序，查找到指定数量的元素
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInOrder 按照时间或分数获取降序的id列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// // 多次请求，造成RTT耗时非常大
	//data = make([]int64, 0, len(ids))
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZsetPrefix + id)
	//	// 查找key中分数是1的元素的数量 -> 统计每篇帖子的赞成票的数量
	//	v := client.ZCount(ctx, key, "1", "1").Val() // ZCount会返回分数在min和max范围内的成员数量
	//	data = append(data, v)
	//}

	// 使用pipeline依次发送多条命令，节省RTT时间
	pipeline := client.Pipeline()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZsetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		// 类型转换
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用zinterstore 把分区的set与按帖子分数的zset取交集 生成一个新的zset
	// 针对新的zset，按之前的逻辑取数据
	orderKey := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZset)
	}
	// 社区的key
	// 例如：community:2143
	cKey := getRedisKey(keysCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存key减少zinterstore执行的次数

	// key是新表的key，orderKey是post:score或者post:time
	// 利用缓存key减少zinterscore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	// 判断key是否存在
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 如果缓存已经过期
	if client.Exists(ctx, key).Val() < 1 {
		// 不存在,需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			// 由于set的score默认是0，
			// 因此如果将zset和set融合并保留zset的score，
			// 应该在zset和set的score之间取最大值
			Aggregate: "MAX",
			Keys:      []string{cKey, orderKey},
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 一分钟后，key会被删除, 本质是缓存
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 存在的话就直接根据key查询ids
	return getIDsFromKey(key, p.Page, p.Size)
}
