package apis

import (
  "os"
  "bytes"
  "net/http"
  "encoding/json"
)

var userUrl string

const privateKey = "hhhhhhhhhhhhhhhhhh"

func init()  {
  if os.Getenv("ADMIN_API_ENV") == "prod" {
    userUrl= "http://47.106.136.96:8069/admin/refresh?"
  } else {
    userUrl= "http://47.106.136.96:8069/admin/refresh?"
  }

}

func Reflash(uid int) error {
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