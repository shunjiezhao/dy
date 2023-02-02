package user

import (
	"context"
	"first/kitex_gen/user"
	"first/pkg/constants"
	"first/service/api/rpc/mock"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

var (
	serverAddr = "http://" + constants.ApiServerAddress
	testUserA  = "DouYinTestUserA"
	testUserB  = "DouYinTestUserB"
)

func getHandler(t *testing.T) (*gin.Engine, *gomock.Controller, *mock.MockRpcProxyIFace) {
	engine := gin.New()
	ctrl := gomock.NewController(t) // 需要去关闭
	face := mock.NewMockRpcProxyIFace(ctrl)
	service := Service{rpc: face}
	InitRouter(engine, &service)

	return engine, ctrl, face
}

func newExpect(t *testing.T, serverAddr string) (*httpexpect.Expect, *gomock.Controller, *mock.MockRpcProxyIFace) {
	handler, controller, face := getHandler(t)
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	}), controller, face
}

//mock 的功能:
//返回user_id
//生成token
//校验生成都token能否得到user_id(将jwt的now函数定格)
func TestRegister(t *testing.T) {
	expect, controller, face := newExpect(t, serverAddr)
	defer controller.Finish()
	ctx := context.Background()
	gomock.InOrder(
		face.EXPECT().Register(ctx, &user.RegisterRequest{
			UserName: testUserA,
			PassWord: testUserA,
		}).Return(int64(-1), nil),
	)
	_ = expect.POST("/douyin/user/register/").
		WithQuery("username", testUserA).WithQuery("password", testUserA).
		WithFormField("username", testUserA).WithFormField("password", testUserA).
		Expect().Status(http.StatusOK).JSON().Object()

	//1. 已经存在的用户 再次注册
	//getTestUserToken(testUserA, expect)
	//2. 注册其他用户 token
	//3. 登陆该账户返回 token
}

func getTestUserToken(user string, e *httpexpect.Expect) (userId int, token string) {
	maxRetry := 2
	for i := 0; i < maxRetry; i++ {
		loginResp := e.GET("/douyin/user/login").WithQuery("username", user).WithQuery("password", user).
			WithFormField("username", user).WithFormField("password", user).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginToken := loginResp.Value("token").String()
		loginToken.Length().Gt(0)
		token = loginToken.Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
		if len(token) == 0 || userId == 0 {
			registerResp := e.POST("/douyin/user/register").
				WithQuery("username", user).WithQuery("password", user).
				WithFormField("username", user).WithFormField("password", user).
				Expect().Status(http.StatusOK).JSON().Object()
			userId = 0
			token = registerResp.Value("token").String().Raw()
			userId = int(registerResp.Value("user_id").Number().Raw())
		} else {
			break
		}
	}
	return userId, token
}
