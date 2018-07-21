package models

import (
	"admin/utils"
	"errors"
	"fmt"
	//google "code.google.com/a_game/src/models"
	"strconv"
)

type UserEx struct {
	BaseModel     `xorm:"-"`
	Uid           int64  `xorm:"not null pk comment(' 用户ID') BIGINT(11)"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)"`
	AffirmTime    int64  `xorm:"comment('实名认证时间') BIGINT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      int64  `xorm:"comment('邀请者id') BIGINT(11)"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)"`
	AffirmCount   int    `xorm:"default 0 comment('实名认证的次数') TINYINT(4)"`
	ChannelName   string `xorm:"not null default '' comment('邀请的渠道名称') VARCHAR(100)" json:"channel_name"`
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
	AffirmTime         int64  `xorm:"comment('实名认证时间') BIGINT(20)"`
	AffirmCount        int    `xorm:"default 0 comment('实名认证的次数') TINYINT(4)"`
	RealNameVerifyMark int    //一级实名认证
	GoogleVerifyMark   int    //google 认证
	TWOVerifyMark      int    //二级认证
	PhoneVerifyMark    int    //电话认证
	EMAILVerifyMark    int    //邮箱认证
	TotalCNY           int64  // 账户的折合总资产
	TotalCurrentCNY    int64  //法币账户折合
	LockCurrentCNY     int64  // 法币折合冻结CNY
	TotalTokenCNY      int64  //币币账户折合
	LockTokenCNY       int64  //bibi 折合冻结CNY

}

type FirstDetail struct {
	UserEx       `xorm:"extends"`
	Account      string `xorm:"comment('账号') unique VARCHAR(64)"`
	SecurityAuth int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	VerifyMark   int    //一级实名认证状态
}

func (f *FirstDetail) TableName() string {
	return "user_ex"
}

func (w *WebUser) TableName() string {
	return "user"
}

func (w *UserGroup) TableName() string {
	return "user"
}

type InviteGroup struct {
	WebUser      `xorm:"extends"`
	Account          string `xorm:"comment('账号') unique VARCHAR(64)"`
	Email            string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Phone            string `xorm:"comment('手机') unique VARCHAR(64)"`
	Status           int    `xorm:"default 0 comment('用户状态，0正常，1冻结') INT(11)"`
	InviteCount int
}

func (w *UserEx) TableName() string {
	return "user_ex"
}

//邀请人列表
func (w*UserEx)GetInviteInfoList(uid,page,rows int ,date uint64,name ,account string) (*ModelList,error){
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("INNER", "user", "user.uid=user_ex.uid")
	query =query.Cols("user.account","user_ex.register_time","user_ex.channel_name")
	query =query.Where("`user_ex`.`invite_id`=?",uid)
	if name!=``{
		temp:=fmt.Sprintf("channer_name=%s",name)
		query = query.Where(temp)
	}
	if account!=``{
		temp:=fmt.Sprintf("user.account=%s",account)
		query =query.Where(temp)
	}
	fmt.Println("刷选时间=",date)
	if date!=0{
		query = query.Where("`user_ex`.`register_time` BETWEEN ? AND ? ", date, date+86400)
	}
	tempQuery:=*query
	count,err:=tempQuery.Count(&UserEx{})
	if err!=nil{
		return nil ,err
	}
	offset,modelList :=w.Paging(page,rows,int(count))
	list:=make([]InviteGroup,0)
	err=query.Limit(modelList.PageSize,offset).Find(&list)
	if err!=nil{
		return nil,err
	}
	modelList.Items =list
	return  modelList,nil
}

//p2-5好友邀请

func (w *WebUser) GetInViteList(page, rows int, search string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user_ex.uid=user.uid and user_ex.invite_id!=''")
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	temp := *query
	count, err := temp.Where("`user_ex`.`invite_id`!=''").Count(&WebUser{})
	if err != nil {
		return nil, err
	}
	fmt.Println("count=",count)

	offset, modelList := w.Paging(page, rows, int(count))
	list := make([]InviteGroup, 0)
	err=query.Limit(modelList.PageSize, offset).Find(&list)
	if err!=nil{
		return nil,err
	}

	sql := "SELECT `user_ex`.`invite_id`, COUNT(*) AS counts FROM `user_ex` WHERE `user_ex`.`invite_id`!='' GROUP BY  `user_ex`.`invite_id` "
	value, err := query.QueryString(sql)
	if err != nil {
		return nil, err
	}

	for _, v := range value {
		// v["counts"]
		uid, _ := strconv.Atoi(v["uid"])
		//if err!=nil{
		//	continue
		//}
		count, _ := strconv.Atoi(v["counts"])
		//if err!=nil{
		//	continue
		//}
		//邀请的人数
		for i,_:=range list{
			if uint64(uid) == list[i].Uid {
				list[i].InviteCount = count
				break
			}
		}

		fmt.Printf("%#v\n", v)
		//for i,vv:=range v {
		//
		//}
	}



	modelList.Items = list
	return modelList, nil

}

//二级认证审核
func (w *WebUser) SecondAffirmLimit(uid, status int) error {
	engine := utils.Engine_common
	query := engine.Where("uid=?", uid)
	temp := *query
	wu := new(WebUser)
	has, err := temp.Get(wu)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("用户不存在！！")
	}
	if status == utils.AUTH_NIL {
		wu.SecurityAuth = wu.SecurityAuth &^ utils.AUTH_TWO
	}
	if status == utils.AUTH_TWO {
		wu.SecurityAuth = wu.SecurityAuth ^ utils.AUTH_TWO // 为实名状态标识
	}
	_, err = query.Update(&WebUser{
		SecurityAuth: wu.SecurityAuth,
	})
	if err != nil {
		return err
	}
	return nil
}

//审核实名
func (w *WebUser) FirstAffirmLimit(uid, status int) error {
	engine := utils.Engine_common
	query := engine.Where("uid=?", uid)
	temp := *query
	wu := new(WebUser)
	has, err := temp.Get(wu)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("用户不存在！！")
	}

	if status == utils.AUTH_NIL {
		wu.SecurityAuth = wu.SecurityAuth &^ utils.AUTH_FIRST
	}
	if status == utils.AUTH_FIRST {
		wu.SecurityAuth = wu.SecurityAuth ^ utils.AUTH_FIRST // 16 为实名状态标识
	}
	_, err = query.Update(&WebUser{
		SecurityAuth: wu.SecurityAuth,
	})
	if err != nil {
		return err
	}

	return nil
}

//单个用户的认证详情
func (w *FirstDetail) GetFirstDetail(uid int) (*FirstDetail, error) {
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("INNER", "user", "user.uid=user_ex.uid")
	query = query.Where("user_ex.uid=?", uid)
	query = query.Cols("user_ex.register_time", "user_ex.uid", "user_ex.real_name", "user_ex.identify_card", "user_ex.affirm_time", "user_ex.affirm_count", "user.account", "user_ex.nick_name", "user.security_auth")
	temp := *query
	has, err := temp.Exist(&FirstDetail{})
	if err != nil {
		fmt.Println("bukexue ")
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在！！")
	}
	u := new(FirstDetail)
	_, err = query.Get(u)
	if err != nil {
		return nil, err
	}
	if u.SecurityAuth&utils.AUTH_FIRST == utils.AUTH_FIRST {
		u.VerifyMark = 1
	}
	return u, nil
}

//拉取一级实名认证列表
func (w *WebUser) GetFirstList(page, rows, status, cstatus int, time uint64, search string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("user.uid")
	query = query.Cols("user_ex.real_name", "user.uid", "user_ex.register_time", "user.phone", "user_ex.nick_name", "user.email", "user.security_auth", "user.status")
	query = query.Join("INNER", "user_ex", "user_ex.uid=user.uid")
	if status != 0 {
		query = query.Where("`user`.`status`=?", status)
	}
	if cstatus == -1 { //未通过
		query = query.Where("security_auth & ? !=?", utils.AUTH_FIRST, utils.AUTH_FIRST)
	}
	if cstatus == utils.AUTH_FIRST {
		query = query.Where("security_auth & ? =?", utils.AUTH_FIRST, utils.AUTH_FIRST)
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

	for i, _ := range list {
		if list[i].SecurityAuth&utils.AUTH_FIRST == utils.AUTH_FIRST {
			list[i].RealNameVerifyMark = 1
		}
	}
	modelList.Items = list
	return modelList, nil
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
		query = query.Where("`user`.`security_auth` & ? =? ", verify, verify)
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
	for _, v := range list {
		if v.SecurityAuth&utils.AUTH_EMAIL == utils.AUTH_EMAIL {
			v.EMAILVerifyMark = 1
		}
		if v.SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			v.TWOVerifyMark = 1
		}
		if v.SecurityAuth&utils.AUTH_FIRST == utils.AUTH_FIRST {
			v.RealNameVerifyMark = 1
		}
		if v.SecurityAuth&utils.AUTH_GOOGLE == utils.AUTH_GOOGLE {
			v.GoogleVerifyMark = 1
		}
		if v.SecurityAuth&utils.AUTH_PHONE == utils.AUTH_PHONE {
			v.PhoneVerifyMark = 1
		}
	}
	modelList.Items = list
	return modelList, nil
}

//获取用户列表
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
	//认证 判段
	for index, v := range list {
		if v.SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			list[index].TWOVerifyMark = 1
		}
	}
	return list, nil
}

func (w *WebUser) GetUserListForUid(uid []uint64) ([]UserGroup, error) {
	engine := utils.Engine_common
	query := engine.Desc("user.uid")
	query = query.Join("INNER", "user_ex", "user.uid=user_ex.uid")
	query = query.In("user.uid", uid)
	list := make([]UserGroup, 0)
	err := query.Find(&list)
	if err != nil {
		return nil, err
	}
	//认证 判段
	for index, v := range list {
		if v.SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			list[index].TWOVerifyMark = 1
		}
	}
	return list, nil
}
