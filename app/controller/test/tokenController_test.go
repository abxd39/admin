package test

import (
	"admin/constant"
	"admin/utils/test"
	"net/http"
	"testing"
)

func Test_Token_Fee_Trend(t *testing.T) {
	c := test.NewAPIClient()

	// 请求参数
	c.AddParam("token_id", "1")
	c.AddParam("date_begin", "2018-08-09")
	c.AddParam("date_end", "2018-09-19")

	// cookie
	cookies := []*http.Cookie{
		LOGIN_COOKIE,
	}

	// 发送
	resp, err := c.Get("token/fee_trend", cookies)
	if err != nil {
		t.Fatal(err)
	}

	// 判断错误码
	if resp.Code != constant.RESPONSE_CODE_OK {
		t.Fatalf(resp.Msg)
	}
}
