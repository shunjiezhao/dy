package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
)

func fatalOnerr(err error) {
	if err != nil {
		if err != redis.Nil {
			panic(err)
		}
	}
}
func TestRedis(t *testing.T) {
	ctx := context.Background()
	var (
		err   error
		me    int64 = 1
		other int64 = 2
	)
	user := NewUserRedis(me, testRedis, ctx)
	user.redis.Del(ctx, GetKey(Follow, me))
	user.redis.Del(ctx, GetKey(Follow, other))
	user.redis.Del(ctx, GetKey(Follower, me))
	user.redis.Del(ctx, GetKey(Follower, other))

	err = user.FollowUser(other)
	fatalOnerr(err)

	// 关注列表是否有
	iFollowHe, err := user.redis.SIsMember(ctx, GetKey(Follow, me), other).Result()
	fatalOnerr(err)
	if iFollowHe == false {
		t.Fatal("没有关注")
	}
	// 粉丝列表是否有
	otherFollowerHasMe, err := user.redis.SIsMember(ctx, GetKey(Follower, other), me).Result()
	fatalOnerr(err)
	if otherFollowerHasMe == false {
		t.Fatal("粉丝没我")
	}
	// 取消关注
	err = user.UnFollowUser(other)
	fatalOnerr(err)
	// 关注列表是否有
	iFollowHe, err = user.redis.SIsMember(ctx, GetKey(Follow, me), other).Result()
	fatalOnerr(err)
	if iFollowHe == true {
		t.Fatal("取消关注操作失败")
	}
	// 粉丝列表是否有
	otherFollowerHasMe, err = user.redis.SIsMember(ctx, GetKey(Follower, other), me).Result()
	fatalOnerr(err)
	if otherFollowerHasMe == true {
		t.Fatal("取消关注操作失败")
	}
}
func TestLikeVideo(t *testing.T) {
	ctx := context.Background()
	var (
		err          error
		likeVideoId1 int64 = 1
		likeVideoId2 int64 = 2
		me           int64 = 2
	)
	testRedis.Del(ctx, GetVideoKey(LikeUser, likeVideoId1))
	testRedis.Del(ctx, GetVideoKey(LikeUser, likeVideoId2))
	testRedis.Del(ctx, GetKey(LikeVideo, me))
	user := NewUserRedis(me, testRedis, ctx)
	// 1.  喜欢视频
	fatalOnerr(user.LikeVideo(likeVideoId1))
	fatalOnerr(user.LikeVideo(likeVideoId2))
	// 2. 查看个人喜欢列表以及, 信息点赞数
	result1, err := testRedis.SMembers(ctx, GetVideoKey(LikeUser, likeVideoId1)).Result()
	assert.Len(t, result1, 1)
	assert.Equal(t, intToString(me), result1[0])

	fatalOnerr(err)
	result2, err := testRedis.SMembers(ctx, GetVideoKey(LikeUser, likeVideoId2)).Result()
	assert.Len(t, result2, 1)
	assert.Equal(t, intToString(me), result2[0])
	fatalOnerr(err)

	userLike, err := testRedis.ZRange(ctx, GetKey(LikeVideo, me), 0, -1).Result()
	assert.Len(t, userLike, 2)
	assert.Equal(t, userLike[0], intToString(likeVideoId1))
	assert.Equal(t, userLike[1], intToString(likeVideoId2))
	fatalOnerr(err)

}

func TestGetVideoInfo(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())
	var (
		videoId  int64 = 1
		palyUrl        = "play_url"
		coverUrl       = "cover_url"
		userId   int64 = 3
		author   int64 = rand.Int63n(100000)
	)

	// 1. 加入信息
	testRedis.Set(ctx, GetVideoKey(PlayUrl, videoId), palyUrl, time.Second)
	testRedis.Set(ctx, GetVideoKey(CoverUrl, videoId), coverUrl, time.Second)
	// 2.
	testRedis.LPush(ctx, GetVideoKey(Comment, videoId), palyUrl)
	testRedis.Expire(ctx, GetVideoKey(Comment, videoId), time.Second)

	testRedis.SAdd(ctx, GetVideoKey(LikeUser, videoId), userId)
	testRedis.Expire(ctx, GetVideoKey(LikeUser, videoId), time.Second)

	testRedis.Set(ctx, GetVideoKey(Author, videoId), author, time.Second)
	testRedis.Set(ctx, GetVideoKey(Author, videoId), author, time.Second)

	result, err := testRedis.Eval(ctx, getVideoInfoLua, []string{
		0: GetVideoKey(PlayUrl, videoId),
		1: GetVideoKey(CoverUrl, videoId),
		2: GetVideoKey(Comment, videoId),
		3: GetVideoKey(LikeUser, videoId),
		4: GetVideoKey(Author, videoId),
	}, userId).Result()

	fatalOnerr(err)
	res, ok := result.([]interface{})
	assert.Equal(t, true, ok)

	assert.Len(t, res, 6)
	assert.Equal(t, palyUrl, res[0].(string))
	assert.Equal(t, coverUrl, res[1].(string))
	assert.Equal(t, int64(1), res[2].(int64))
	assert.Equal(t, int64(1), res[3].(int64))
	assert.Equal(t, int64(1), res[4].(int64))
	assert.Equal(t, author, res[5].(int64))
}

func TestGetFeedsInfo(t *testing.T) {
	ctx := context.Background()
	var (
		videoId  int64 = 1
		palyUrl        = "play_url"
		coverUrl       = "cover_ur"
		userId   int64 = 3
		author   int64 = 1000
	)
	// 1. 加入信息

	var i int64
	for i = 2; i >= 0; i-- {
		testRedis.LPush(ctx, videoFeedListKey, videoId+i)
		testRedis.Expire(ctx, videoFeedListKey, time.Second)

		testRedis.Set(ctx, GetVideoKey(PlayUrl, videoId+i), palyUrl+intToString(i), time.Second)
		testRedis.Set(ctx, GetVideoKey(CoverUrl, videoId+i), coverUrl+intToString(i), time.Second)

		testRedis.LPush(ctx, GetVideoKey(Comment, videoId+i), palyUrl+intToString(i))
		testRedis.Expire(ctx, GetVideoKey(Comment, videoId+i), time.Second)

		testRedis.SAdd(ctx, GetVideoKey(LikeUser, videoId+i), userId)
		testRedis.Expire(ctx, GetVideoKey(LikeUser, videoId+i), time.Second)

		testRedis.Set(ctx, GetVideoKey(Author, videoId+i), author, time.Second)
		testRedis.Set(ctx, GetVideoKey(Author, videoId+i), author, time.Second)
	}

	result, err := testRedis.Eval(ctx, getFeedsInfoLua, []string{
		0: videoFeedListKey,
		1: videoPlayUrlPrefix,
		2: videoCoverUrlPrefix,
		3: videoCommentPrefix,
		4: videoLikeUserPrefix,
		5: videoAuthorPrefix,
	}, 0, 2, userId).Result()
	fatalOnerr(err)

	res, ok := result.([]interface{})
	if !ok {
		t.Fatal("失败")
	}

	assert.Len(t, res, 3)
	for j, tval := range res {
		val, ok := tval.([]interface{})
		if !ok {
			t.Fatal("失败")
		}
		i = int64(j)
		assert.Len(t, val, 7)
		assert.Equal(t, palyUrl+intToString(i), val[0].(string))
		assert.Equal(t, coverUrl+intToString(i), val[1].(string))
		assert.Equal(t, int64(1), val[2].(int64))
		assert.Equal(t, int64(1), val[3].(int64))
		assert.Equal(t, int64(1), val[4].(int64))
		assert.Equal(t, videoId+i, val[5].(int64))
		assert.Equal(t, author, val[6].(int64))
	}

}

func TestGetUserInfo(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var (
		me         int64 = 1
		userId     []int64
		userCounts int64 = 3
		base             = rand.Int63n(100000)
		name             = "username."
		counts           = map[int64]bool{}

		args []interface{}
	)

	ctx := context.Background()

	var i int64
	for i = 0; i < userCounts; i++ {
		uId := base + i
		// 1. 用户姓名
		testRedis.Set(ctx, GetKey(Name, uId), name+intToString(uId), time.Second*10)

		testRedis.SAdd(ctx, GetKey(Follow, uId), me)
		testRedis.Expire(ctx, GetKey(Follow, uId), time.Second*10)

		testRedis.SAdd(ctx, GetKey(Follower, uId), me)
		testRedis.Expire(ctx, GetKey(Follower, uId), time.Second*10)
		counts[uId] = true

		userId = append(userId, uId)
		args = append(args, uId)

	}
	user := NewUserRedis(me, testRedis, context.Background())
	eval := user.redis.Eval(user.ctx, getUserInfoLua, []string{
		0: userNamePrefix,
		1: userFollowSetPrefix,
		2: userFollowerSetPrefix,
		3: intToString(user.userId),
	}, args...)
	result, err := eval.Result()
	fatalOnerr(err)
	res, ok := result.([]interface{})
	if !ok {
		t.Fatal("失败")
	}
	assert.Len(t, res, int(userCounts))
	for i := 0; i < len(res); i++ {
		val := res[i].([]interface{})
		if assert.Len(t, val, 5) == false {
			t.Fatal("失败")
		}
		j := val[idIdx].(int64)
		assert.Equal(t, counts[j], true)
		counts[j] = false
		assert.Equal(t, val[nameIdx].(string), name+intToString(j))
		assert.Equal(t, val[isFollowIdx].(int64), int64(1))
		assert.Equal(t, val[followCountIdx].(int64), int64(1))
		assert.Equal(t, val[followerCountIdx].(int64), int64(1))
	}
}

var testRedis *redis.Client

func TestGetFollowList(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())
	var (
		me         int64 = 1
		userId     []int64
		userCounts int64 = 3
		base             = rand.Int63n(100000)
		name             = "username."
		counts           = map[int64]bool{}
		i          int64
		expire     = time.Second
	)

	for i = 0; i < userCounts; i++ {
		uId := base + i
		// 1. 用户姓名
		testRedis.Set(ctx, GetKey(Name, uId), name+intToString(uId), expire)

		testRedis.SAdd(ctx, GetKey(Follow, me), uId)
		testRedis.Expire(ctx, GetKey(Follow, me), expire)
		counts[uId] = true

		userId = append(userId, uId)

	}
	user := NewUserRedis(me, testRedis, ctx)
	list, err := user.FollowUserList()
	fatalOnerr(err)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})
	if assert.Len(t, list, int(userCounts)) == false {
		t.Fatal("失败")
	}
	for i = 0; int(i) < len(list); i++ {
		assert.Equal(t, list[i].Id, base+i)
		assert.Equal(t, list[i].Name, name+intToString(base+i))

	}

}
func TestMain(t *testing.M) {
	testRedis = InitRedis()
	os.Exit(t.Run())
}
