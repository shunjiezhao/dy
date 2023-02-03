package user

import (
	"first/pkg/constants"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

var (
	serverAddr = "http://" + constants.ApiServerAddress
	testUserA  = "DouYinTestUserA"
	testUserB  = "DouYinTestUserB"
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
