package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type VendorApi struct{}

var userUrl = ""
var awardUrl = ""
var privateKey = "hhhhhhhhhhhhhhhhhh"

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

func InitAwardUrl(url, key string) {
	if key != `` {
		privateKey = key
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
	fmt.Println(userUrl+"/admin/refresh?")
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
	body,err:=ioutil.ReadAll(result.Body)
	if err!=nil{
		return err
	}
	rsp:=&struct {
		Code int
		Msg string
	}{}
	err=json.Unmarshal(body,rsp)
	if err!=nil{
		return err
	}
	fmt.Println(rsp)
	if rsp.Code!=0{
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
		return errors.New("failed!!!")
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
	//url :=userHost
	request, err := http.NewRequest("POST", walletUrl+"/wallet/sendrawtx?", reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}
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
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	rsp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return err
	}
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
	CnyPriceInt int    `json:"cny_price_int"`
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
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", walletUrl+"market/cny_prices?", reader)
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
	fmt.Println("1111111111111111111111111")
	fmt.Println(userUrl + url)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		fmt.Println("123111111")
		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("1472555556565656")
		return nil, err
	}
	fmt.Println("6666666666666666666666666")
	result := new(Token)
	fmt.Println(string(body))
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
	fmt.Println("result.Data.list", result.Data.List)
	return result.Data.List, nil
}
