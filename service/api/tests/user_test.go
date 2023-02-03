package user

import (
	"testing"
)

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
