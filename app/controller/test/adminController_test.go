package test

import (
	"net/http"
	"testing"

	"admin/constant"
	"admin/utils/test"
)

func Test_Admin_List(t *testing.T) {
	c := test.NewAPIClient()

	// 请求参数
	c.AddParam("page", "1")
	c.AddParam("rows", "10")

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Get("admin/list", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}

	t.Log(resp.Data)
}

func Test_Admin_Add(t *testing.T) {
	c := test.NewAPIClient()

	// 请求参数
	c.AddParam("name", TEST_ADMIN_NAME)
	c.AddParam("nickname", TEST_ADMIN_NAME)
	c.AddParam("pwd", "123456")
	c.AddParam("re_pwd", "123456")
	c.AddParam("role_ids", "1,2,3")

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Post("admin/add", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}
}
