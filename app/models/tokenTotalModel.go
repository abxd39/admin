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
	UserToken `xorm:"extends"`
	NickName  string
	Phone     string
	Email     string
	Btc       float32 `xorm:"-"` //折合比特币总数
	Balance   float64 `xorm:"-"` // 这和人民币总数
	Status    int     //账号状态
}

// var Total []PersonalProperty

// func NewTotal() {
// 	Total = make([]PersonalProperty, 0)
// }

func (t *PersonalProperty) TableName() string {
	return "user_token"
}
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
func (t *PersonalProperty) TotalUserBalance(page, rows, status int, search string) (*ModelList, error) {
	//查 用户表
	engine := utils.Engine_token
	query := engine.Alias("ut")
	query = query.Join("LEFT", "g_common.user u", "u.uid=ut.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=ut.uid")
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(ut.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)

	}

	countQuery:=*query
	count,err:=countQuery.Count(&PersonalProperty{})
	if err!=nil{
		return nil,err
	}
	 offset,mList:=t.Paging(page,rows,int(count))
	 list:=make([]PersonalProperty,0)
	 err=query.Limit(mList.PageSize,offset).Find(&list)
	if err!=nil{
		return nil,err
	}
	mList.Items = list
	//if status != 0 || search != `` {
	//	list, err := new(WebUser).GetAllUser(page, rows, status, search)
	//	if err != nil {
	//		return nil, err
	//	}
	//	var uid []int64
	//	userlist, Ok := list.Items.([]UserGroup)
	//	if !Ok {
	//		return nil, errors.New("assert failed!!")
	//	}
	//	for _, v := range userlist {
	//		uid = append(uid, v.Uid)
	//	}
	//	fmt.Printf("TotalUserBalance%#v\n", uid)
	//
	//	query := engine.Desc("uid")
	//	query = query.In("uid", uid)
	//	tempQuery := *query
	//	count, err := tempQuery.Count(&UserToken{})
	//	if err != nil {
	//		return nil, err
	//	}
	//	offset, modelList := t.Paging(page, rows, int(count))
	//	tokenlist := make([]PersonalProperty, 0)
	//	err = query.Limit(modelList.PageSize, offset).Find(&tokenlist)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	for index, _ := range tokenlist {
	//
	//		for _, ob := range userlist {
	//			if tokenlist[index].Uid == uint64(ob.Uid) {
	//				tokenlist[index].Email = ob.Email
	//				tokenlist[index].NickName = ob.NickName
	//				tokenlist[index].Phone = ob.Phone
	//
	//			}
	//
	//		}
	//	}
	//
	//	modelList.Items = tokenlist
	//	return modelList, nil
	//}
	////去重找uid的所有uid
	//query := engine.Desc("uid")
	//tempQuery := *query
	//
	//count, err := tempQuery.Count(&PersonalProperty{})
	//if err != nil {
	//	return nil, err
	//}
	//offset, modelList := t.Paging(page, rows, int(count))
	//tokenlist := make([]PersonalProperty, 0)
	//
	//query = query.Limit(modelList.PageSize, offset)
	//countQuery := *query
	//err = query.Find(&tokenlist)
	//if err != nil {
	//	return nil, err
	//}
	//countList := make([]UserToken, 0)
	//err = countQuery.Distinct("uid").Find(&countList)
	//if err != nil {
	//	return nil, err
	//}
	////根据uid 获取用户资料
	//uidlist := make([]uint64, 0)
	//for _, v := range countList {
	//	uidlist = append(uidlist, v.Uid)
	//}
	//
	//ulist, err := new(UserGroup).GetUserListForUid(uidlist)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for i, _ := range tokenlist {
	//	for _, value := range ulist {
	//		if tokenlist[i].Uid == uint64(value.Uid) {
	//			tokenlist[i].Phone = value.Phone
	//			tokenlist[i].NickName = value.NickName
	//			tokenlist[i].Email = value.Email
	//		}
	//	}
	//}
	//modelList.Items = tokenlist
	//return modelList, nil
	return mList, nil
}
