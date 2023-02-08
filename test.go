package main

import (
	"context"
	redis2 "first/pkg/redis"
	"first/service/api/handlers"
	"fmt"
	"github.com/u2takey/go-utils/json"
)

type User struct {
	*handlers.User
}

var testlua = `for i=1,2,3 do
    print(1)
end `

func main() {
	cli := redis2.InitRedis()
	eval := cli.Eval(context.Background(), testlua, nil)
	fmt.Println(eval.Result())

}
func mainw() {
	U := &User{&handlers.User{
		Id:            1,
		Name:          "打发水电费 ",
		FollowCount:   1,
		FollowerCount: 2,
		IsFollow:      false,
	}}
	cli := redis2.InitRedis()
	ctx := context.Background()
	push := cli.LPush(ctx, "test.key.value", U)
	fmt.Println(push)
	pop := cli.LPop(ctx, "test.key.value")
	bytes, err := pop.Bytes()
	if err != nil {
		panic(err)
	}
	fmt.Println(json.Unmarshal(bytes, U))
	fmt.Println(U.User)
}
