package user

import (
	"context"
	"first/kitex_gen/user"
	"first/pkg/constants"
	"first/pkg/errno"
	"first/pkg/middleware"
	user2 "first/service/api/handlers/common/user"
	"first/service/api/rpc/mock"
	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	serverAddr = "http://" + constants.ApiServerAddress
	testUserA  = "DouYinTestUserA"
	testUserB  = "DouYinTestUserB"
)

// TODO: 样例过滤 SQL 注入 [非法字符]

func getHandler(t *testing.T) (*gin.Engine, *gomock.Controller, *mock.MockRpcProxyIFace) {
	engine := gin.New()
	ctrl := gomock.NewController(t) // 需要去关闭
	face := mock.NewMockRpcProxyIFace(ctrl)
	service := user2.Service{rpc: face}
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

//检查生成token
func TestRegister(t *testing.T) {
	expect, controller, face := newExpect(t, serverAddr)
	defer controller.Finish()
	ctx := context.Background()
	assert := assert.New(t)

	gomock.InOrder(
		face.EXPECT().Register(ctx, &user.RegisterRequest{UserName: testUserA, PassWord: testUserA}).Return(int64(0),
			errno.ServiceErr),
		face.EXPECT().Register(ctx, &user.RegisterRequest{UserName: testUserA, PassWord: testUserA}).Return(int64(1), nil),
	)
	tests := []struct {
		name      string
		userName  string
		shouldErr bool
		err       error

		userId int64 // token
	}{
		{"缺少用户名/密码", "", true, errno.ParamErr, -1},
		{"注册失败", testUserA, true, errno.ServiceErr, -1},
		{"注册成功", testUserA, false, errno.Success, 1},
	}

	for _, test := range tests {
		resp := expect.POST("/douyin/user/register/").
			WithQuery("username", test.userName).WithQuery("password", test.userName).
			WithFormField("username", test.userName).WithFormField("password", test.userName).
			Expect().Status(http.StatusOK).JSON().Object()
		assert.NotNilf(resp, "%s:http resp is nil", test.name)

		// 1.检查 code
		var respCode int64
		if test.shouldErr {
			respCode = errno.ConvertErr(test.err).ErrCode
		}
		resp.ContainsKey("status_code").ValueEqual("status_code", respCode)
		// 2. 检查msg
		if test.shouldErr == false {
			validateToken(resp, assert, test.name, test.userId)
		}

	}

	//1. 已经存在的用户 再次注册
	//getTestUserToken(testUserA, expect)
	//2. 注册其他用户 token
	//3. 登陆该账户返回 token
}

func validateToken(resp *httpexpect.Object, assert *assert.Assertions, name string,
	wantUserId int64) {
	jwt, _ := middleware.JwtMiddle()
	jwt.TokenLookup = "header: Authorization"
	// 3. 如果成功 检查返回的token
	// len(token) > 0
	token := resp.Value("token").String()
	token.Length().Gt(0)

	claim, err := jwt.ParseTokenString(token.Raw())
	assert.Nilf(err, "%s: jwt parse token error: %v", name, err) // err == nil
	assert.NotNilf(claim, "%s: jwt get claim nil", name)         // claim != nil

	fromToken := jwt2.ExtractClaimsFromToken(claim)          // 验证是本人的token 嘛
	assert.NotNilf(fromToken, "%s: jwt get claim nil", name) // claim != nil
	id := fromToken[constants.IdentityKey].(float64)
	assert.Equalf(int64(id), wantUserId, "%s: token claim user_id is not equal", name)
}

func TestLogin(t *testing.T) {
	expect, controller, face := newExpect(t, serverAddr)
	defer controller.Finish()
	ctx := context.Background()
	assert := assert.New(t)

	gomock.InOrder(
		face.EXPECT().CheckUser(ctx, &user.CheckUserRequest{UserName: testUserB, PassWord: testUserB}).Return(int64(0),
			errno.AuthorizationFailedErr),
		face.EXPECT().CheckUser(ctx, &user.CheckUserRequest{UserName: testUserA, PassWord: testUserA}).Return(int64(1), nil),
	)
	tests := []struct {
		name      string
		userName  string
		shouldErr bool
		err       error

		userId int64 // token
	}{
		// 检验参数的 样例
		{"缺少用户名/密码", "", true, errno.ParamErr, -1},

		// 参数检验成功
		{"登陆失败", testUserB, true, errno.AuthorizationFailedErr, -1},
		{"登陆成功", testUserA, false, errno.Success, 1},
	}

	for _, test := range tests {
		resp := expect.POST("/douyin/user/login/").
			WithQuery("username", test.userName).WithQuery("password", test.userName).
			WithFormField("username", test.userName).WithFormField("password", test.userName).
			Expect().Status(http.StatusOK).JSON().Object()
		if resp == nil {
			t.Fatalf("%s:http resp is nil", test.name)
		}

		// 1.检查 code
		var respCode int64
		if test.shouldErr {
			respCode = errno.ConvertErr(test.err).ErrCode
		}
		resp.ContainsKey("status_code").ValueEqual("status_code", respCode)
		// 2. 检查msg
		if test.shouldErr == false {
			validateToken(resp, assert, test.name, test.userId)
		}

	}

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
