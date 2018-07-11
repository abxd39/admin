package models

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"
)

var accessKeyId string = "LTAIcJgRedhxruPq"
var accessKeySecret string = "d7p6tWRfy0B2QaRXk7q4mb5seLROtb"
var host string = "sdun.oss-cn-hongkong.aliyuncs.com"
var expire_time int64 = 3600
var upload_dir string = "sdun/"

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func (p *PolicyToken) get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func (p *PolicyToken) Get_policy_token() PolicyToken {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = p.get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = upload_dir
	policyToken.Policy = string(debyte)
	// response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}
	// return string(response)
	return policyToken
}

// func hello(w http.ResponseWriter, r *http.Request) {
// 	response := get_policy_token()
// 	w.Header().Set("Access-Control-Allow-Methods", "POST")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	io.WriteString(w, response)
// }

// func main() {
// 	http.HandleFunc("/", hello)
// 	http.ListenAndServe(":1234", nil)
// }
