package user

import (
	"first/pkg/constants"
	"first/service/api/router"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

var (
	serverAddr = "http://" + constants.ApiServerAddress
	testUserA  = "douyinTestUserA"
	testUserB  = "douyinTestUserB"
)

func getHandler() *gin.Engine {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	//TODO: nil 代替为实现 use proxy rpc 的接口(mock)
	router.InitRouter(engine, nil)
	return engine
}

func newExpect(t *testing.T, serverAddr string) *httpexpect.Expect {

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(getHandler()),
			Jar:       httpexpect.NewCookieJar(),
		},
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

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

func TestMain(t *testing.M) {

}
