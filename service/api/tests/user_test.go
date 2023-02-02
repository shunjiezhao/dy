package user

import (
	"first/pkg/constants"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func newExpect(t *testing.T, serverAddr string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

var (
	serverAddr = "http://" + constants.ApiServerAddress
	testUserA  = "douyinTestUserA"
	testUserB  = "douyinTestUserB"
)

//TODO: 利用 Mock 将rpc 调用打包
//mock 的功能:
//返回user_id
//生成token
//校验生成都token能否得到user_id(将jwt的now函数定格)
func TestRegister(t *testing.T) {
	//1. 已经存在的用户 再次注册
	getTestUserToken(testUserA, newExpect(t, serverAddr))
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
