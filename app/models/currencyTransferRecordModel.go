package models

type CurrencyTransferRecordModel struct {
	Id         int64  `xorm:"id pk" json:"id"'`
	Uid        int32  `xorm:"uid" json:"uid"`
	TokenId    int32  `xorm:"token_id" json:"token_id"`
	TokenName  string `xorm:"token_name" json:"token_name"`
	Num        int64  `xorm:"num" json:"num"`
	States     int32  `xorm:'states' json:"states"`
	CreateTime int64  `xorm:'create_time' json:"create_time"`
}

func (CurrencyTransferRecordModel) TableName() string {
	return "transfer_record" // g_currency库
}
