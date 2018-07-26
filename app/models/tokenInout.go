package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
)

//冲提币明细流水表
type TokenInout struct {
	BaseModel   `xorm:"-"`
	Id          int    `xorm:"not null pk autoincr comment('自增id') INT(11)"`
	Uid         int64  `xorm:"not null comment('用户id') INT(11)"`
	Opt         int    `xorm:"not null comment('操作方向 1 充币 2 提币') TINYINT(4)"`
	Txhash      string `xorm:"not null comment('交易hash') VARCHAR(200)"`
	From        string `xorm:"not null comment('打款方') VARCHAR(42)"`
	To          string `xorm:"not null comment('收款方') VARCHAR(42)"`
	Amount      int64  `xorm:"not null comment('金额') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('提币手续费') BIGINT(20)"`
	Value       string `xorm:"not null comment('原始16进制转账数据') VARCHAR(32)"`
	Chainid     int    `xorm:"not null comment('链id') INT(11)"`
	Contract    string `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Tokenid     int    `xorm:"not null comment('币种id') INT(11)"`
	States      int    `xorm:"not null comment('充提币状态 1正在充提币 对应操作状态为 撤销，2 已完成，3已取消') TINYINT(1)"`
	TokenName   string `xorm:"not null comment('币种名称') VARCHAR(10)"`
	CreatedTime string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间 提币创建时间') TIMESTAMP"`
	DoneTime    string `xorm:"not null default '0000-00-00 00:00:00' comment('充币到账时间') TIMESTAMP"`
	Remarks     string `xorm:"not null comment('备注信息') VARCHAR(100)"`
}

type TokenInoutGroup struct {
	TokenInout `xorm:"extends"`
	NickName   string `xorm:"-" json:"nick_name"`
	Phone      string `xorm:"-" json:"phone"`
	Email      string `xorm:"-" json:"email"`
	Status     int    `xorm:"-" json:"status"`
}

func (t *TokenInoutGroup) TableName() string {
	return "token_inout"
}

//提币 充币 p3-1-0 充币 提币管理
func (t *TokenInout) GetTokenInList(page, rows, uStatus, status, tokenId, opt int, search, date string) (*ModelList, error) {
	engine := utils.Engine_wallet
	query := engine.Desc("states")
	//两个方向  用户信息库和 钱包库
	if uStatus != 0 || search != `` {
		mList, err := new(WebUser).GetAllUser(page, rows, uStatus, search)
		if err != nil {
			return nil, err
		}
		value, ok := mList.Items.([]UserGroup)
		if !ok {
			return nil, errors.New("assert []webUser type failed!!")
		}
		uidList := make([]int64, 0)
		for _, v := range value {
			uidList = append(uidList, v.Uid)
		}
		if len(uidList) < 1 {
			//没有匹配刷选条件的用户
			return nil, nil
		}
		query = query.In("uid", uidList)
		countQuery := *query
		count, err := countQuery.Count(&TokenInout{})
		offset, modelList := t.Paging(page, rows, int(count))
		list := make([]TokenInoutGroup, 0)
		err = query.Limit(modelList.PageSize, offset).Find(&list)
		if err != nil {
			return nil, err
		}
		//
		for i, _ := range list {
			for _, v := range value {
				if list[i].Uid == v.Uid {
					list[i].NickName = v.NickName
					list[i].Phone = v.Phone
					list[i].Email = v.Email
					list[i].Status = v.Status
					break
				}
			}
		}
		modelList.Items = list
		return modelList, nil
	} else {
		if tokenId != 0 {
			query = query.Where("token_id=?", tokenId)
		}
		if status != 0 {
			query = query.Where("states=?", status)
		}
		if opt != 0 {
			query = query.Where("opt=?", opt)
		}
		if date != `` {
			subst := date[:11] + "23:59:59"
			fmt.Println(subst)
			sql := fmt.Sprintf("create_time  BETWEEN '%s' AND '%s' ", date, subst)
			query = query.Where(sql)
		}
		countQuery := *query
		count, err := countQuery.Count(&TokenInout{})
		if err != nil {
			return nil, err
		}
		offset, modelList := t.Paging(page, rows, int(count))
		list := make([]TokenInoutGroup, 0)
		err = query.Limit(modelList.PageSize, offset).Find(&list)
		if err != nil {
			return nil, err
		}
		uidList := make([]uint64, 0)
		for _, v := range list {
			uidList = append(uidList, uint64(v.Uid))
		}
		if len(uidList) < 1 {
			//没有匹配刷选条件的用户
			return nil, nil
		}
		uList, err := new(UserGroup).GetUserListForUid(uidList)
		if err != nil {
			return nil, err
		}
		for i, _ := range list {
			for _, v := range uList {
				if v.Uid == list[i].Uid {
					list[i].TokenName = v.NickName
					list[i].Phone = v.Phone
					list[i].Email = v.Email
					list[i].Status = v.Status
					break
				}
			}
		}
		modelList.Items = list
		return modelList, nil
	} //else
	return nil, nil
}

//提币管理
func (t *TokenInout) OptTakeToken(id, uid int) error {
	engine := utils.Engine_wallet
	query := engine.Desc("id")
	//t:=new(TokenInout)
	has, err := query.Where("id=? and uid=?", id, uid).Get(t)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("rescind failed !!")
	}

	_, err = query.Where("id=? and uid=?", id, uid).Update(&TokenInout{
		States: 3,
	})
	if err != nil {
		return err
	}
	return nil

}
