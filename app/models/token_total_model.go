package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

type UserToken struct {
	BaseModel `xorm:"-"`
	Id        int64
	Uid       uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId   int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance   int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen    int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version   int    `xorm:"version"`
}

//资产
type PersonalProperty struct {
	Uid      int
	NickName string
	Phone    string
	Email    string
	Btc      float32 //折合比特币总数
	Balance  float64 // 这和人民币总数
	Status   int     //账号状态
	token    []UserToken
}

// var Total []PersonalProperty

// func NewTotal() {
// 	Total = make([]PersonalProperty, 0)
// }

func (u *UserToken) GetTokenDetailOfUid(uid, token_id int) ([]UserToken, error) {
	if uid < 0 {
		return nil, errors.New("uid is illegal")
	}
	engine := utils.Engine_token
	list := make([]UserToken, 0)
	err := engine.Where("uid=?", uid).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//所有用户 的全部币币资产
//第一步get 所有用户
func (t *PersonalProperty) TotalUserBalance(page, rows, status int) (*ModelList, error) {
	//查 用户表

	list, err := new(WebUser).GetAllUser(page, rows, status)
	if err != nil {
		return nil, err
	}
	var uid []uint64
	userlist, Ok := list.Items.([]UserGroup)
	if !Ok {
		return nil, errors.New("assert failed!!")
	}
	for _, v := range userlist {
		uid = append(uid, v.Uid)
	}
	fmt.Printf("TotalUserBalance%#v\n", uid)
	engine := utils.Engine_token
	token := make([]UserToken, 0)
	err = engine.In("uid", uid).Find(&token)
	if err != nil {
		return nil, err
	}
	Total := make([]PersonalProperty, 0)
	for _, ob := range userlist {
		pp := &PersonalProperty{}
		pp.Uid = int(ob.Uid)
		pp.NickName = ob.NickName
		pp.Phone = ob.Phone
		pp.Email = ob.Email
		pp.Status = ob.Status
		for _, result := range token {
			if ob.Uid == result.Uid {
				pp.token = append(pp.token, result)
			}
		}

		Total = append(Total, *pp)
	}
	list.Items = Total
	return list, nil
}
