package models

import (
	"admin/utils"
)

// type Tokens struct {
// 	Id        int    `xorm:"not null pk autoincr INT(11)"`
// 	Name      string `xorm:"default '' comment('货币名称') VARCHAR(64)"`
// 	Detail    string `xorm:"default '' comment('详情地址') VARCHAR(255)"`
// 	Signature string `xorm:"default '' comment('签名方式(eth,eip155,btc)') VARCHAR(255)"`
// 	Chainid   int    `xorm:"default 0 comment('链id') INT(11)"`
// 	Github    string `xorm:"default '' comment('项目地址') VARCHAR(255)"`
// 	Web       string `xorm:"default '' VARCHAR(255)"`
// 	Mark      string `xorm:"comment('英文标识') CHAR(10)"`
// 	Logo      string `xorm:"default '' comment('货币logo') VARCHAR(255)"`
// 	Contract  string `xorm:"default '' comment('合约地址') VARCHAR(255)"`
// 	Node      string `xorm:"default '' comment('节点地址') VARCHAR(100)"`
// 	Decimal   int    `xorm:"not null default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)"`
// }

type Tokens struct {
	Id     int    `xorm:"not null pk autoincr INT(10)"`
	Name   string `xorm:"not null default '' comment('虚拟货币名称') VARCHAR(20)"`
	CnName string `xorm:"not null default '' comment('虚拟货币中文名称') VARCHAR(20)"`
}

func (t *Tokens) TableName() string {
	return "tokens"
}

//获取数字货币id及名称
func (t *Tokens) GetTokenList() ([]Tokens, error) {
	engine := utils.Engine_currency
	list := make([]Tokens, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//获取数据货币的名称
// func (t *Tokens) GetTokensList() ([]Tokens, error) {
// 	fmt.Println("1321321313212")
// 	engin := utils.Engine_common
// 	tokens := make([]Tokens, 0)
// 	err := engin.Find(&tokens)
// 	if err != nil {
// 		utils.AdminLog.Errorln(err.Error())
// 		return nil, err
// 	}
// 	return tokens, nil
// }
