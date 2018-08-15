package models

import "admin/utils"
type UserLoginTerminalType struct {
	BaseModel    `xorm:"-"`
	Id           int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	TerminalType int    `xorm:"not null comment('终端类型') TINYINT(4)"`
	TerminalName string `xorm:"not null comment('终端的名称 例如:pc') VARCHAR(100)"`
}

func (u *UserLoginTerminalType) GetTerminalTypeList() (*ModelList, error) {
	engine := utils.Engine_common
	list := make([]UserLoginTerminalType, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	ml := new(ModelList)
	ml.Items = list
	return ml, nil
}
