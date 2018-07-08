package models

import (
	"admin/utils"
	"fmt"
	//google "code.google.com/a_game/src/models"
)

type UserEx struct {
	Uid           int64  `xorm:"not null pk comment(' 用户ID') BIGINT(11)"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      int64  `xorm:"comment('邀请者id') BIGINT(11)"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)"`
}

type WebUser struct {
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
}

type UserGroup struct {
	WebUser            `xorm:"extends"`
	NickName           string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	RegisterTime       int64  `xorm:"comment('注册时间') BIGINT(20)"`
	RealNameVerifyMark int
	GoogleVerifyMark   int
	TWOVerifyMark      int
}

func (w *WebUser) TableName() string {
	return "user"
}

func (w *UserGroup) TableName() string {
	return "user"
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

func (w *WebUser) GetAllUser(page, rows, status int) ([]UserGroup, int, int, error) {
	engine := utils.Engine_common
	if page <= 1 {
		page = 1
	}
	var begin int

	if rows <= 1 {
		rows = 100
	}

	if page > 1 {
		begin = (page - 1) * rows
	} else {
		begin = 0
	}
	users := make([]UserGroup, 0)
	u := new(WebUser)
	count, err := engine.Where("status=?", status).Count(u)
	if err != nil {
		return nil, 0, 0, err
	}
	var total int
	if int(count) > rows {
		total = int(count) / rows
		v := int(count) % rows
		if v >= 0 {
			total = total + 1
		}
	} else {
		total = 1
	}
	err = engine.Join("INNER", "user_ex", "user_ex.uid=user.uid").Limit(rows, begin).Find(&users)
	// sql := fmt.Sprintf("select a.*,b.nick_name,b.register_time from `user` a left join user_ex b on a.uid=b.uid ")
	// fmt.Println("GetAllUser", rows, begin)
	// err = engine.Sql(sql).Limit(rows, begin).Find(&users)
	if err != nil {
		return nil, 0, 0, err
	}
	return users, total, total * rows, nil
}

func (w *WebUser) UserList(page, rows, verify, status int, uname, phone, email string, date, uid int64) ([]*UserGroup, int, int, error) {
	if rows == 0 {
		rows = 50
	}
	engine := utils.Engine_common
	var begin int
	var total int
	var count int64
	var err error
	if page <= 1 {
		begin = 0
	} else {
		begin = rows * (page - 1)
	}

	//数据查询
	users := make([]*UserGroup, 0)
	fmt.Println("vvvvvvvvvvvvvvvvvvv")
	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user.uid=user_ex.uid")
	//筛选条件为 uid
	if uid != 0 {
		query = query.Where("uid=?", uid)
		fmt.Println("uid")
	}
	//刷选条件为电话号码
	if len(phone) != 0 {
		query = query.Where("phone=?", phone)
		fmt.Println("phone")
	}
	//刷选选条件为用户名
	if len(uname) != 0 {
		query = query.Where("status=", status)
		fmt.Println("uname")
	}
	//刷选条件为Email
	if len(email) != 0 {
		query = query.Where("email=", email)
		fmt.Println("email")
	}
	//刷选条件为 用户状态
	if status == 1 {
		query = query.Where("status=", status)
	}
	//刷选条件为用户注册的日期
	if date != 0 {
		query = query.Where("register_time>?", date)
	}
	//刷选条件为用户的验证方式
	if verify != 0 {
		query = query.Where("verify=?", verify)
	}
	//无条件刷选

	fmt.Printf("无条件刷选%#v\n", rows)
	tempQuery := query
	query = query.Limit(rows, begin)
	err = query.Find(&users)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, 0, 0, err
	}
	count, err = tempQuery.Count(&WebUser{})
	if err != nil {
		return nil, 0, 0, err
	}
	total = int(count) / rows
	return users, total, total * rows, nil
}
