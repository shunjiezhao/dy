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
	prefix     = "/douyin/user"
	user       = "test-user"
)

func TestRegister(t *testing.T) {
	getTestUserToken(user, newExpect(t, serverAddr))
}

func getTestUserToken(user string, e *httpexpect.Expect) (int, string) {
	registerResp := e.POST(prefix+"/register").
		WithQuery("username", user).WithQuery("password", user).
		WithFormField("username", user).WithFormField("password", user).
		Expect().Status(http.StatusOK).JSON().Object()
	userId := 0
	token := registerResp.Value("token").String().Raw()
	if len(token) == 0 {
		loginResp := e.GET(prefix+"/login").WithQuery("usrename", user).WithQuery("password", user).
			WithFormField("username", user).WithFormField("password", user).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginToken := loginResp.Value("token").String()
		loginToken.Length().Gt(0)
		token = loginToken.Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
	} else {
		userId = int(registerResp.Value("user_id").Number().Raw())
	}
	return userId, token
}
