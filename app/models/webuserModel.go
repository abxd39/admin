package models

import (
	"admin/utils"

	"admin/apis"
	"errors"
	"fmt"

	"time"
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
	Uid              int64  `xorm:"not null pk autoincr comment('用户ID') BIGINT(11)"`
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
	Status           int    `xorm:"default 1 comment('用户状态，1正常，2冻结') INT(11)"`
	SecurityAuth     int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	WhiteList        int    `xorm:"not null default 2 comment('用户白名单 1为白名单 免除交易手续费，2 需要缴纳交易手续费') TINYINT(4)"`
	SetTardeMark     int    `xorm:"default 0 comment('交易密码是否设置状态') INT(8)"`
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
	UserEx      `xorm:"extends"`
	Account     string `xorm:"comment('账号') unique VARCHAR(64)"`
	Email       string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Phone       string `xorm:"comment('手机') unique VARCHAR(64)"`
	Status      int    `xorm:"default 0 comment('用户状态，0正常，1冻结') INT(11)"`
	InviteCount int
}

func (w *UserEx) TableName() string {
	return "user_ex"
}

//邀请人统计表—账号：18888888888
func (w *UserEx) GetInviteInfoList(uid, page, rows int, date uint64, name, account string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("LEFT", "user", "user.uid=user_ex.uid")
	query = query.Cols("user_ex.uid", "user.account", "user_ex.register_time", "user_ex.channel_name")
	query = query.Where("`user_ex`.`invite_id`=?", uid)
	if name != `` {
		temp := fmt.Sprintf("channer_name=%s", name)
		query = query.Where(temp)
	}
	if account != `` {
		temp := fmt.Sprintf("user.account=%s", account)
		query = query.Where(temp)
	}
	fmt.Println("刷选时间=", date)
	if date != 0 {
		query = query.Where("`user_ex`.`register_time` BETWEEN ? AND ? ", date, date+86400)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&UserEx{})
	if err != nil {
		return nil, err
	}
	offset, modelList := w.Paging(page, rows, int(count))
	list := make([]InviteGroup, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

//p2-5好友邀请 ___ 有邀请用户注册的列表

func (w *UserEx) GetInViteList(page, rows int, search string) (*ModelList, error) {
	engine := utils.Engine_common
	countSql := "SELECT COUNT(*) cnt FROM (%s) newtable"
	sql := "SELECT u.uid,u.`email`,u.`phone`,ue.`real_name`,ue.nick_name,u.account,ue.`invite_id`,u.`status`,tmp.cnt invite_count FROM (SELECT invite_id,COUNT(invite_id) cnt FROM user_ex GROUP BY invite_id) tmp JOIN user_ex ue ON ue.uid=tmp.invite_id JOIN `user` u ON u.uid=tmp.invite_id WHERE tmp.invite_id!=0"
	limitSql := " LIMIT %d OFFSET %d "

	if search != `` {
		temp := fmt.Sprintf(" AND CONCAT(IFNULL(`ue`.`uid`,''),IFNULL(`u`.`phone`,''),IFNULL(`ue`.`nick_name`,''),IFNULL(`u`.`email`,'')) LIKE '%%%s%%' ", search)
		sql += temp
	}
	//query := engine.SQL(contentSql + sql)

	type test struct {
		Cnt int
	}
	var temp test
	_, err := engine.SQL(fmt.Sprintf(countSql, sql)).Get(&temp)

	fmt.Println("cnt=", temp.Cnt)
	fmt.Println("page=", page, "rows=", rows, "cnt=", temp.Cnt)
	offset, modelList := w.Paging(page, rows, temp.Cnt)
	list := make([]InviteGroup, 0)
	limitSql = fmt.Sprintf(limitSql, modelList.PageSize, offset)
	fmt.Println("sql2=", sql+limitSql)
	err = engine.SQL(sql + limitSql).Find(&list)
	if err != nil {
		return nil, err
	}

	//fmt.Println("resultList=",list)
	modelList.Items = list
	return modelList, nil

}

//二级认证审核
func (w *WebUser) SecondAffirmLimit(uid, status int) error {
	//err :=new(apis.VendorApi).AddAwardToken(uid)
	//if err!=nil{
	//	//sess.Rollback()
	//	fmt.Println(err.Error())
	//	fmt.Println("赠送奖励失败")
	//}
	//return nil
	engine := utils.Engine_common
	sess := engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	wu := new(WebUser)
	has, err := sess.Where("uid=?", uid).Get(wu)
	if err != nil {
		sess.Rollback()
		return err
	}
	if !has {
		sess.Rollback()
		return errors.New("用户不存在！！")
	}
	us := new(UserSecondaryCertification)
	has, err = sess.Table("user_secondary_certification").Where("uid=?", uid).Get(us)
	if err != nil {
		return err
	}
	if !has {
		sess.Rollback()
		return errors.New("not exists!!")
	}
	ReverseSidePath := us.ReverseSidePath
	InHandPicturePath := us.InHandPicturePath
	PositivePath := us.PositivePath
	if status == utils.AUTH_NIL {
		//审核不通过删除数据
		//oss
		wu.SecurityAuth = wu.SecurityAuth &^ utils.AUTH_TWO
		wu.SetTardeMark = wu.SetTardeMark ^ utils.APPLY_FOR_SECOND_NOT_ALREADY //二级认证没有通过
		wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_SECOND            //申请撤销
		fmt.Println("uid=", uid, "------------------>二级认证", wu.SetTardeMark)
		if _, err = sess.Table("user_secondary_certification").Where("uid=?", uid).AllCols().Update(&UserSecondaryCertification{
			ReverseSidePath:       "",
			InHandPicturePath:     "",
			PositivePath:          "",
			VerifyTime:            0,
			VideoRecordingDigital: "",
		}); err != nil {
			sess.Rollback()
			return err
		}
		a := new(Article)
		if ReverseSidePath != `` {
			a.DeletFileToAliCloud(ReverseSidePath)
		}
		if InHandPicturePath != `` {
			a.DeletFileToAliCloud(InHandPicturePath)
		}
		if PositivePath != `` {
			a.DeletFileToAliCloud(PositivePath)
		}
		_, err = sess.Where("uid=?", uid).Cols("security_auth", "set_tarde_mark").Update(&WebUser{
			SecurityAuth: wu.SecurityAuth,
			SetTardeMark: wu.SetTardeMark,
		})
		if err != nil {
			sess.Rollback()
			return err
		}
		sess.Commit()
		err = new(apis.VendorApi).Reflash(uid)
		if err != nil {
			fmt.Println("缓存清理失败!!!")
		}
		return nil
	}
	if status == utils.AUTH_TWO {
		wu.SecurityAuth = wu.SecurityAuth ^ utils.AUTH_TWO // 为实名状态标识
	}
	//审核过之后不管通没通过审核 实名申请的状态一律设为为 未申请状态
	wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_SECOND
	wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_SECOND_NOT_ALREADY //二级认证没有通过
	//wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_SECOND            //申请撤销
	_, err = sess.Where("uid=?", uid).Cols("security_auth", "set_tarde_mark").Update(&WebUser{
		SecurityAuth: wu.SecurityAuth,
		SetTardeMark: wu.SetTardeMark,
	})
	if err != nil {
		sess.Rollback()
		return err
	}
	err =new(apis.VendorApi).AddAwardToken(uid)
	if err!=nil{
		sess.Rollback()
		fmt.Println("赠送奖励失败")
	}
	sess.Commit()
	err = new(apis.VendorApi).Reflash(uid)
	if err != nil {
		fmt.Println("缓存清理失败!!!")
	}

	return nil
}

//审核实名
func (w *WebUser) FirstAffirmLimit(uid, status int) error {

	engine := utils.Engine_common
	sess := engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	//temp := *query
	wu := new(WebUser)
	has, err := sess.Where("uid=?", uid).Get(wu)
	if err != nil {
		sess.Rollback()
		return err
	}
	if !has {
		sess.Rollback()
		return errors.New("用户不存在！！")
	}

	//
	if status == utils.AUTH_NIL {
		//审核不过删除信息
		has, err := sess.Table("user_ex").Where("uid=?", uid).Exist(&UserEx{})
		if err != nil {
			sess.Rollback()
			return err
		}
		if !has {
			return errors.New("not exists")
		}
		if _, err = sess.Table("user_ex").Where("uid=?", uid).Cols("identify_card", "real_name").Update(&UserEx{
			IdentifyCard: "",
			RealName:     "",
		}); err != nil {
			sess.Rollback()
			return err
		}
		wu.SecurityAuth = wu.SecurityAuth &^ utils.AUTH_FIRST
		wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_FIRST            //撤销申请
		wu.SetTardeMark = wu.SetTardeMark ^ utils.APPLY_FOR_FIRST_NOT_ALREADY //没有通过
	}
	if status == utils.AUTH_FIRST {
		wu.SecurityAuth = wu.SecurityAuth ^ utils.AUTH_FIRST // 16 为实名状态标识
		wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_FIRST
		wu.SetTardeMark = wu.SetTardeMark &^ utils.APPLY_FOR_FIRST_NOT_ALREADY //没有通过
	}
	//删除 申请状态

	_, err = sess.Where("uid=?", uid).Cols("security_auth", "set_tarde_mark").Update(&WebUser{
		SecurityAuth: wu.SecurityAuth,
		SetTardeMark: wu.SetTardeMark,
	})
	if err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	err =new(apis.VendorApi).Reflash(uid)
	if err != nil {
		fmt.Println("缓存清理失败!!!")
	}

	return nil
}

//单个用户的认证详情
func (w *FirstDetail) GetFirstDetail(uid int) (*FirstDetail, error) {
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("INNER", "user", "user.uid=user_ex.uid")
	query = query.Where("user_ex.uid=?", uid)
	//query = query.Cols("user_ex.register_time", "user_ex.uid", "user_ex.real_name", "user_ex.identify_card", "user_ex.affirm_time", "user_ex.affirm_count", "user.account", "user_ex.nick_name", "user.security_auth")
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
	//query = query.Cols("user_ex.real_name", "user.uid", "user_ex.register_time", "user.phone", "user_ex.nick_name", "user.email", "user.security_auth", "user.status")
	query = query.Join("INNER", "user_ex", "user_ex.uid=user.uid")
	query = query.Where("user.set_tarde_mark &2 =2 or user.set_tarde_mark&16=16 or user.security_auth&16=16")
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

func (w *WebUser) GetTotalUser() (int, int, int, error) {
	engine := utils.Engine_common
	Count := &struct {
		TotalCount  int
		TodayCount  int
		UpDayCount  int
		UpWeekCount int
	}{}

	_, err := engine.SQL("select count(*) total_count from g_common.user_ex").Get(Count)
	if err != nil {
		return 0, 0, 0, err
	}
	//获取当前时间
	current := time.Now().Unix() //当前时间戳
	currentDay := time.Now().Format("2006-01-02 15:04:05")
	//ExampleParseInLocation()
	//叫上日 涨幅
	//24*3600 =
	currentDay = fmt.Sprintf("u.days='%s'", currentDay[:10])
	_, err = engine.SQL("SELECT COUNT(*) today_count FROM (SELECT FROM_UNIXTIME(register_time,'%Y-%m-%d')days FROM g_common.user_ex ) u WHERE " + currentDay).Get(Count)
	upDayUnix := current - 86400
	tm := time.Unix(upDayUnix, 0)
	upDaySql := tm.Format("2006-01-02 15:04:05")
	upDaySql = fmt.Sprintf("u.days='%s'", upDaySql[:10])
	_, err = engine.SQL("SELECT COUNT(*) up_day_count FROM (SELECT FROM_UNIXTIME(register_time,'%Y-%m-%d')days FROM g_common.user_ex ) u WHERE " + upDaySql).Get(Count)
	if err != nil {
		return 0, 0, 0, err
	}
	//叫上周同日 涨幅
	upWeekUnix := current - 86400*7
	tw := time.Unix(upWeekUnix, 0)
	upWeekStr := tw.Format("2006-01-02 15:04:05")
	upWeekStr = fmt.Sprintf("u.days='%s'", upWeekStr[:10])
	_, err = engine.SQL("SELECT COUNT(*) up_week_count FROM (SELECT FROM_UNIXTIME(register_time,'%Y-%m-%d')days FROM g_common.user_ex ) u WHERE " + upWeekStr).Get(Count)
	if err != nil {
		return 0, 0, 0, err
	}
	upDay := Count.TodayCount - Count.UpDayCount

	upWeek := Count.TodayCount - Count.UpWeekCount

	return Count.TotalCount, upDay, upWeek, nil
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
func (w *WebUser) GetCurrencyList(uid []int64, verify int, search string) ([]UserGroup, error) {
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

type TotalProperty struct {
	BaseModel `xorm:"-"`
	Uid int64
	Phone     string
	NickName  string
	Status    int
	Rt       int64 //用户注册时间
	Email     string
	TokenId   uint32
	Ucb       int64 //法币总资产
	Ucf       int64 //法币冻结
	Utb       int64 //币币总资产
	Utc       int64 //币币冻结
}

func (*TotalProperty) TableName() string {
	return "user"
}


func (w *TotalProperty) GetTotalProperty(page, rows, status int, date uint64, search string) (*ModelList, error) {
	engine := utils.Engine_common

	countSql:="SELECT COUNT(u.`uid`) num "
	contentSql:=" SELECT `u`.`uid`, `u`.`email`, `u`.`phone`, `ex`.`nick_name`, `ex`.`register_time` `rt`, `uc`.`balance`  `ucb`, `uc`.`freeze` `ucf`, `ut`.`balance` `utb`, `ut`.`frozen` `utf` "
	mid:="FROM `user` AS `u` LEFT JOIN g_common.user_ex ex ON u.uid=ex.uid LEFT JOIN g_currency.user_currency uc ON u.uid=uc.uid LEFT JOIN g_token.user_token ut ON u.uid=ut.uid "
	sql:=fmt.Sprintf(" WHERE (ex.register_time BETWEEN %d AND %d)",date,date+86400)
	if status != 0 {
		appendSql := fmt.Sprintf("and  u.status=%d ", status)
		sql+=appendSql
	}
	if search != `` {
		appendSql:= fmt.Sprintf(" and concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		sql += appendSql
	}
	Count:=&struct{
	 Num int
	}{}
	_,err := engine.SQL(countSql+mid+sql).Get(Count)
	if err != nil {
		return nil, err
	}
	offset, mList := w.Paging(page, rows, int(Count.Num))
	list := make([]TotalProperty, 0)
	limitSql:=fmt.Sprintf("limit %d offset %d",mList.PageSize, offset)
	err = engine.SQL(contentSql+mid+sql+limitSql).Find(&list)
	if err != nil {
		return nil, err
	}
	mList.Items = list
	return mList, nil
}
