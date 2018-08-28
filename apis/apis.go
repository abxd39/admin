package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type VendorApi struct{}

var userUrl = ""
var awardUrl = ""
var privateKey = "hhhhhhhhhhhhhhhhhh"
var verifykey = "32zBKHYCjK8ZBWbwCG1HarNZqOBBbKLodsCI1V20"

var walletUrl = ""

func InitUserUrl(remoteUrl, localUrl, key string) {
	if key != `` {
		privateKey = key
	}
	if os.Getenv("ADMIN_API_ENV") == "prod" {
		userUrl = remoteUrl + userUrl
		walletUrl = remoteUrl

	} else {
		userUrl = localUrl + userUrl
		walletUrl = localUrl
	}

}

func InitAwardUrl(url, key, verify string) {
	if key != `` {
		privateKey = key
	}
	if verify != `` {
		verifykey = verify
	}
	awardUrl = url
}

func (VendorApi) Reflash(uid int) error {
	//fmt.Println(userUrl)
	params := make(map[string]interface{})
	params["uid"] = uid
	params["key"] = privateKey
	fmt.Println(params)
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	fmt.Println(userUrl + "/admin/refresh?")
	request, err := http.NewRequest("POST", userUrl+"/admin/refresh?", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int
		Msg  string
	}{}
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	fmt.Println(rsp)
	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}
	return nil
}

//后台审核通过之后 赠送平台币
func (VendorApi) AddAwardToken(uid int) error {
	fmt.Println(awardUrl)
	params := make(map[string]interface{})
	params["uid"] = uid
	params["key"] = privateKey
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	request, err := http.NewRequest("POST", awardUrl+"/admin/register_reward?", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	//mapCode :=make(map[interface{}]interface{})

	// _=result
	type ReturnValue struct {
		Code int `json:"code"`
		//Msg string `json:"msg"`
		//Data string `json:"data"`
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	//fmt.Println("000000000000000000000000000000000_result.Body=》",string(body))
	var returnValue ReturnValue
	err = json.Unmarshal(body, &returnValue)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", returnValue)
	if returnValue.Code != 0 {
		return errors.New("赠送货币失败！！")
	}
	return nil
}

//提币审核 获取签名
//wallet/signtx
func (VendorApi) GetTradeSigntx(uid, tid int, addr, mount string) (string, error) {
	params := make(map[string]interface{})
	params["uid"] = uid
	params["token_id"] = tid
	params["to"] = addr
	params["amount"] = mount
	params["gasprice"] = 100
	paramByts, err := json.Marshal(params)
	reader := bytes.NewReader(paramByts)
	//url :=userHost
	fmt.Println(params)
	fmt.Println(walletUrl + "/wallet/signtx")
	request, err := http.NewRequest("POST", walletUrl+"/wallet/signtx?", reader)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {

		return "", err
	}
	rsp := &struct {
		Code int               `json:"code"`
		Data map[string]string `json:"data"`
		Msg  string            `json:"msg"`
	}{}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(walletUrl + "/wallet/signtx?")
	fmt.Println(string(body))
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return "", err
	}
	if rsp.Code != 0 {
		return "", errors.New(rsp.Msg)
	}
	fmt.Println(rsp)
	v, ok := rsp.Data["signtx"]
	if !ok {
		return "", errors.New("返回值错误")
	}
	return v, nil

}

// 后台审核通过之后发送提币请求 eth
func (VendorApi) PostOutToken(uid, tid, id int, sign string) error {
	params := make(map[string]interface{})
	params["uid"] = uid
	params["token_id"] = tid
	params["apply_id"] = id
	params["signtx"] = sign
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost                                           wallet/sendrawtx
	request, err := http.NewRequest("POST", walletUrl+"/wallet/sendrawtx", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Do(request)

	//if err != nil {
	//	fmt.Println(err)
	//}
	//rsp := &struct {
	//	Code int    `json:"code"`
	//	Msg  string `json:"msg"`
	//}{}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(walletUrl+"/wallet/sendrawtx")
	//fmt.Println("返回值----->",string(body))
	//err = json.Unmarshal(body, rsp)
	//if err != nil {
	//	return err
	//}
	//if rsp.Code != 0 {
	//	return errors.New(rsp.Msg)
	//}
	return nil
}

//后台审核通过之后发送提币请求 btc
//wallet/biti_btc
func (VendorApi) PostOutTokenBtc(uid, tid, id int, addr, mount string) error {
	params := make(map[string]interface{})
	params["uid"] = uid
	params["token_id"] = tid
	params["apply_id"] = id
	params["address"] = addr
	params["amount"] = mount
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	request, err := http.NewRequest("POST", walletUrl+"/wallet/biti_btc?", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Do(request)
	//if err != nil {
	//	return err
	//}
	//rsp := &struct {
	//	Code int    `json:"code"`
	//	Msg  string `json:"msg"`
	//}{}
	//body, err := ioutil.ReadAll(result.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(walletUrl + "/wallet/biti_btc?")
	//fmt.Println(string(body))
	//err = json.Unmarshal(body, rsp)
	//if err != nil {
	//	return err
	//}
	//if rsp.Code != 0 {
	//	return errors.New(rsp.Msg)
	//}
	return nil
}

//提币审核撤销
func (VendorApi) RevokeOutToken(uid, tid, num int64) error {
	params := make(map[string]interface{})
	params["uid"] = uid
	params["ukey"] = fmt.Sprintf("%d_%d", time.Now().UnixNano(), uid)
	params["key"] = verifykey
	params["token_id"] = tid
	params["num"] = num
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	fmt.Println(userUrl + "/wallet/cancel_fronze")
	request, err := http.NewRequest("POST", userUrl+"/wallet/cancel_fronze", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int
		Msg  string
	}{}
	fmt.Println(userUrl + "/wallet/cancel_fronze")
	fmt.Println(string(body))
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	fmt.Println(rsp)
	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}
	return nil
}

//币种和汇率
//人民币价格获取
//market/cny_prices
type Price struct {
	CnyPrice    string `json:"cny_price"`
	UsdPrice    string `json:"usd_price"`
	TokenId     int    `json:"token_id"`
	CnyPriceInt int64  `json:"cny_price_int"`
	UsdPriceInt int64  `json:"usd_price_int"`
}
type temp struct {
	List []Price `json:"list"`
}

func (VendorApi) GetTokenCnyPriceList(tid []int) ([]Price, error) {
	params := make(map[string]interface{})
	params["token_id"] = tid
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(userUrl + "/market/cny_prices?")
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", userUrl+"/market/cny_prices?", reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(result.Body)
	type base struct {
		Code int    `json:"code"`
		Data temp   `json:"data"`
		Msg  string `json:"msg"`
	}
	b := new(base)
	err = json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}
	if b.Code != 0 {
		return nil, errors.New(b.Msg)
	}
	if len(b.Data.List) <= 0 {
		return nil, errors.New("Get Price failed !!!")
	}
	return b.Data.List, nil
}

type Cny struct {
	Uid        uint64 `json:"uid"`
	BalanceCny string `json:"balance_cny"`
	FrozenCny  string `json:"frozen_cny"`
	TotalCny   string `json:"total_cny"`
	BalanceCnyInt int64 `json:"balance_cny_int"`
	FrozenCnyInt int64 `json:"frozen_cny_int"`
}

type listCny struct {
	List []Cny `json:"list"`
}

type Token struct {
	Code int     `json:"code"`
	Data listCny `json:"data"`
	Msg  string  `json:"msg"`
}

//获取币币账户 余额折合 冻结折合
//uid
//mark bibi 法币的标识
//admin/users_total
func (VendorApi) GetCny(uid []uint64, mark int) ([]Cny, error) {
	params := make(map[string]interface{})
	url := ""
	if mark == 1 { //bibi 账户资产
		url = "/admin/users_total" //admin/users_total
		params["uid"] = uid
		params["key"] = privateKey
	} else if mark == 2 { //法币账户资产
		url = "/admin/get_users_balances"
		params["uids"] = uid
		params["key"] = privateKey
	} else {
		return nil, errors.New("unknown")
	}
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", userUrl+url, reader)
	if err != nil {
		return nil, err
	}
	fmt.Println(userUrl + url)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	rsp, err := client.Do(request)
	if err != nil {

		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {

		return nil, err
	}

	result := new(Token)
	//fmt.Println(string(body))
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if result.Code != 0 {

		return nil, errors.New(result.Msg)
	}
	//if len(result.Data.list)<=0{
	//	return nil, errors.New("return value empty")
	//}
	//fmt.Println("result.Data.list", result.Data.List)
	return result.Data.List, nil
}

//法币后台审核
//admin/currency_order_confirm
func (VendorApi) CurrencyVerityPass(id int64) error {
	params := make(map[string]interface{})
	params["id"] = id
	params["key"] = privateKey
	fmt.Println(params)
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	fmt.Println(userUrl + "/admin/currency_order_confirm")
	request, err := http.NewRequest("POST", userUrl+"/admin/currency_order_confirm", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int
		Msg  string
	}{}
	fmt.Println(string(body))
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	fmt.Println(rsp)
	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}
	return nil
}

//法币后台审核撤销
//admin/currency_order_cancel
func (VendorApi) CurrencyRevoke(id int64) error {
	params := make(map[string]interface{})
	params["id"] = id
	params["key"] = privateKey
	fmt.Println(params)
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	fmt.Println(userUrl + "/admin/currency_order_cancel")
	request, err := http.NewRequest("POST", userUrl+"/admin/currency_order_cancel", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int
		Msg  string
	}{}
	fmt.Println(string(body))
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	fmt.Println(rsp)
	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}
	return nil
}

//func (v VendorApi)Test() {
//	signtx := v.httpPost2()
//	result := v.httpPost(signtx)
//	fmt.Println(result)
//}

//func (VendorApi)httpPost(signtx string) string {
//	resp, err := http.Post("http://47.106.136.96:8069/wallet/sendrawtx",
//		"application/x-www-form-urlencoded",
//		strings.NewReader("uid=1&token_id=4&apply_id=7&signtx="+signtx))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(string(body))
//
//	return gjson.Get(string(body),"data.result").String()
//}
//
//
//
//func (VendorApi)httpPost2() string {
//	resp, err := http.Post("http://47.106.136.96:8069/wallet/signtx",
//		"application/x-www-form-urlencoded",
//		strings.NewReader("to=0x870F49783e9d8c9707a72B252a0e56d3b7628F31&gasprice=1&amount=0.0013&uid=10057&token_id=4"))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//
//	return gjson.Get(string(body),"data.signtx").String()
//}
