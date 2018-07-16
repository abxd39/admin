package models

import (
	"admin/utils"
	"errors"
	"fmt"
	//google "code.google.com/a_game/src/models"
)

type UserEx struct {
	BaseModel     `xorm:"-"`
	Uid           int64  `xorm:"not null pk comment(' 用户ID') BIGINT(11)" json:"uid"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)" json:"nick_name"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)" json:"head_sculpture"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)" json:"register_time"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)" json:"invite_code"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)" json:"real_name"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)" json:"identify_card"`
	InviteId      int64  `xorm:"comment('邀请者id') BIGINT(11)" json:"invite_id"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)" json:"invites"`
}

type WebUser struct {
	BaseModel        `xorm:"-"`
	Uid              uint64 `xorm:"not null pk autoincr comment('用户ID') BIGINT(11)"`
	Account          string `xorm:"comment('账号') unique VARCHAR(64)"`
	Pwd              string `xorm:"comment('密码') VARCHAR(255)"`
	Country          string `xorm:"comment('地区号') VARCHAR(32)"`
	Phone            string `xorm:"comment('手机') unique VARCHAR(64)"`
	PhoneVerifyTime  int    `xorm:"comment('手机验证时间') INT(11)"`
	Email            string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	EmailVerifyTime  int    `xorm:"comment('邮箱验证时间') INT(11)"`
	GoogleVerifyId   string `xorm:"comment('谷歌私钥') VARCHAR(128)"`
	GoogleVerifyTime int    `xorm:"comment('谷歌验证时间') INT(255)"`
	SmsTip           int    `xorm:"default 0 comment('短信提醒') TINYINT(1)"`
	PayPwd           string `xorm:"comment('支付密码') VARCHAR(255)"`
	NeedPwd          int    `xorm:"comment('免密设置1开启0关闭') TINYINT(1)"`
	NeedPwdTime      int    `xorm:"comment('免密周期') INT(11)"`
	Status           int    `xorm:"default 0 comment('用户状态，0正常，1冻结') INT(11)"`
	SecurityAuth     int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	WhiteList        int    `xorm:"not null default 2 comment('用户白名单 1为白名单 免除交易手续费，2 需要缴纳交易手续费') TINYINT(4)"`
}

type UserGroup struct {
	WebUser            `xorm:"extends"`
	NickName           string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	RegisterTime       int64  `xorm:"comment('注册时间') BIGINT(20)"`
	RealName           string `xorm:"comment(' 真名') VARCHAR(32)"`
	RealNameVerifyMark int
	GoogleVerifyMark   int
	TWOVerifyMark      int
	Total              int64 // 账户的折合总资产
}

func (w *WebUser) TableName() string {
	return "user"
}

func (w *UserGroup) TableName() string {
	return "user"
}

//获取白名单列表
func (w *WebUser) GetWhiteList(page, rows, status int, time uint64, search string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user_ex.uid=user.uid")
	query = query.Where("`user`.`white_list`=1")
	if status != 0 {
		fmt.Println("status=", status)
		query = query.Where("`user`.`status`=?", status)
		//sql += subsql
	}
	if time != 0 {
		//query = query.Where("")
		query = query.Where("`user_ex`.`register_time` BETWEEN ? AND ? ", time, time+86400)
	}
	if len(search) > 0 {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&WebUser{})
	if err != nil {
		return nil, err
	}
	offset, modelList := w.Paging(page, rows, int(count))

	list := make([]UserGroup, 0)

	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}

	modelList.Items = list
	return modelList, nil
}

//修改黑白名单状态
func (w *WebUser) ModifyWhiteStatus(uid, status int) error {
	engine := utils.Engine_common
	has, err := engine.Where("uid=?", uid).Exist(&WebUser{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("user not exist!!")
	}
	_, err = engine.Where("uid=?", uid).Update(&WebUser{
		WhiteList: status,
	})
	if err != nil {
		return err
	}
	return nil
}

//冻结解冻用户
func (w *WebUser) ModifyUserStatus(uid, status int) error {
	engine := utils.Engine_common
	has, err := engine.Where("uid=?", uid).Exist(&WebUser{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("user not exist!!")
	}
	_, err = engine.Where("uid=?", uid).Update(&WebUser{
		Status: status,
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *WebUser) GetTotalUser() (int, error) {
	engine := utils.Engine_common
	u := new(WebUser)
	count, err := engine.Count(u)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (w *WebUser) GetAllUser(page, rows, status int, search string) (*ModelList, error) {
	engine := utils.Engine_common

	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user_ex.uid=user.uid")
	if status != 0 {
		query = query.Where("status=?", status)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	tempquery := *query
	count, err := tempquery.Count(&WebUser{})
	if err != nil {
		return nil, err
	}
	offset, modelList := w.Paging(page, rows, int(count))
	users := make([]UserGroup, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&users)
	if err != nil {
		return nil, err
	}

	modelList.Items = users
	return modelList, nil
}

func (w *WebUser) UserList(page, rows, verify, status int, search string, date int64) (*ModelList, error) {

	engine := utils.Engine_common

	//数据查询

	query := engine.Desc("user.uid")
	//sql := "FROM `user` a,`user_ex` b WHERE a.uid=b.uid  " //AND b.nick_name LIKE '%w%'

	//query = query.Where(sql)
	query = query.Join("INNER", "user_ex", "user.uid=user_ex.uid")
	//筛选条件为 uid
	var temp string
	if len(search) > 0 {
		temp = fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}

	//刷选条件为 用户状态
	if status != 0 {
		fmt.Println("status=", status)
		query = query.Where("`user`.`status`=?", status)
		//sql += subsql
	}
	//刷选条件为用户注册的日期
	if date != 0 {
		query = query.Where("`user_ex`.`register_time` BETWEEN ? AND ? ", date, date+86400)

	}
	//刷选条件为用户的验证方式
	if verify != 0 {
		//subsql := fmt.Sprintf("AND a.security_auth=%d ", verify)
		query = query.Where("`user`.`security_auth`=? ", verify)
		//sql += subsql
	}
	//无条件刷选
	tempQuery := *query
	count, err := tempQuery.Count(&WebUser{})
	if err != nil {
		return nil, err
	}
	offset, modelList := w.Paging(page, rows, int(count))

	list := make([]UserGroup, 0)

	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}

	modelList.Items = list
	return modelList, nil
}

//
func (w *WebUser) GetCurreryList(uid []uint64, verify int, search string) ([]UserGroup, error) {
	engine := utils.Engine_common

	//数据查询

	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user.uid=user_ex.uid")
	query = query.In("user.uid", uid)

	if verify != 0 {
		query = query.Where("security_auth=?", verify)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	list := make([]UserGroup, 0)
	err := query.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
