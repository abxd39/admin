package models

type TokenChainInout struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	Txhash      string `xorm:"comment('交易hash') VARCHAR(255)"`
	From        string `xorm:"comment('打款地址') VARCHAR(42)"`
	To          string `xorm:"comment('付款地址') VARCHAR(42)"`
	Value       string `xorm:"comment('金额') VARCHAR(30)"`
	Contract    string `xorm:"comment('合约地址') VARCHAR(42)"`
	Chainid     int    `xorm:"comment('链id') INT(11)"`
	Type        int    `xorm:"not null comment('平台转出:1,充值:2') INT(11)"`
	Signtx      string `xorm:"comment('平台转出记录交易签名') VARCHAR(1024)"`
	Tokenid     int    `xorm:"not null comment('币种id') INT(11)"`
	TokenName   string `xorm:"not null comment('币名称') VARCHAR(10)"`
	Uid         int    `xorm:"not null comment('用户id') INT(11)"`
	CreatedTime string `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	InOut       string `xorm:"not null comment('充提币方式 例如：二维码，哈希值') VARCHAR(20)"`
}
