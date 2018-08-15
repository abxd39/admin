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

var userUrl = "/admin/refresh?"
var awardUrl = "/admin/register_reward?"
var awardKey = "hhhhhhhhhhhhhhhhhh"
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
		awardKey = key
	}
	awardUrl = url + awardUrl
}

func (VendorApi) Reflash(uid int) error {
	fmt.Println(userUrl)
	params := make(map[string]interface{})
	params["uid"] = uid
	params["key"] = privateKey
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	request, err := http.NewRequest("POST", userUrl, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}

//后台审核通过之后 赠送平台币
func (VendorApi) AddAwardToken(uid int) error {
	fmt.Println(awardUrl)
	params := make(map[string]interface{})
	params["uid"] = uid
	params["key"] = awardKey
	bytesData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	//url :=userHost
	request, err := http.NewRequest("POST", awardUrl, reader)
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

//提币审核
//wallet/signtx
func (VendorApi) GetTradeSigntx(uid, tid int, addr, mount string) (string, error) {
	fmt.Println("xxxxxxxxxxxxxxxxxxx1")
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

	err = json.Unmarshal(body, rsp)
	if err != nil {
		return "", err
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxx2")
	if rsp.Code != 0 {
		return "", errors.New(rsp.Msg)
	}
	v, ok := rsp.Data["signtx"]
	if !ok {
		return "", errors.New("返回值错误")
	}
	fmt.Printf(v)
	return v, nil

}
