package backstage

import (
	"admin/errors"
	"fmt"

	"admin/app/models"
	"admin/utils"
	"strconv"
	"strings"
	"time"
)

// 管理员s
type User struct {
	models.BaseModel `xorm:"-"`
	Uid              int    `xorm:"not null pk autoincr INT(11)" json:"uid"`
	Name             string `xorm:"not null comment('用户名') VARCHAR(20)" json:"name"`
	NickName         string `xorm:"not null default '' comment('昵称') VARCHAR(60)" json:"nick_name"`
	Pwd              string `xorm:"not null comment('用户登录密码') CHAR(32)" json:"-"`
	Salt             string `xorm:"not null comment('密码加密') CHAR(5)" json:"-"`
	States           int    `xorm:"not null default 1 comment('1正常  0 锁定') TINYINT(4)" json:"states"`
	Remark           string `xorm:"not null default '' comment('备注') VARCHAR(36)" json:"remark"`
	CreateTime       int64  `xorm:"not null comment('创建时间') INT(11)" json:"create_time"`
	UpdateTime       int64  `xorm:"not null comment('修改时间') INT(11)" json:"update_time"`
	LastLoginTime    int64  `xorm:"not null default 0 comment('上次登录时间') INT(11)" json:"last_login_time"`
}

// 登录
func (u *User) Login(name, pwd string) (string, int, error) {
	engine := utils.Engine_backstage
	fmt.Println("login")
	user := &User{}
	_, err := engine.Where("name=?", name).Get(user)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		fmt.Println("login", err.Error())
		return "", 0, err
	}
	//find 如果不存在 数据库是否会返回 一个错误给我
	if user.States == 0 {
		return "", 0, errors.New("该用户已锁定")
	}
	fmt.Printf("数据库中的has值为%s", user.Pwd)
	if utils.Md5(utils.Md5(pwd)+user.Salt) != user.Pwd {
		return "", 0, errors.New("密码不对！！")
	}
	//
	return user.NickName, user.Uid, nil
}

// 管理员列表
func (u *User) List(pageIndex, pageSize int, filter map[string]string) (modelList *models.ModelList, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("uid")

	// 筛选
	query.Where("1=1")
	if v, ok := filter["phone"]; ok {
		query.And("phone like '%?%'", v)
	}

	tempQuery := *query
	count, err := tempQuery.Count(&User{})
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := u.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	var list []User
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	modelList.Items = list

	return
}

// 管理员详情
func (u *User) Get(uid int) (user *User, err error) {
	engine := utils.Engine_backstage
	user = new(User)
	has, err := engine.ID(uid).Get(user)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("管理员不存在或已被删除")
	}

	return
}

// 新增管理员
func (u *User) Add(user *User, roleIds string) (uid int, err error) {
	// 整理数据
	salt := utils.NewRandomString(5)
	now := time.Now().Unix()

	user.Pwd = utils.Md5(utils.Md5(user.Pwd) + salt) // md5两次，第二次带上salt
	user.Salt = salt
	user.CreateTime = now
	user.UpdateTime = now

	// 判断管理员名称是否已存在
	engine := utils.Engine_backstage
	has, err := engine.Where("name=?", user.Name).Get(new(User))
	if err != nil {
		return 0, errors.NewSys(err)
	}
	if has {
		return 0, errors.NewNormal("管理员名称已存在")
	}

	// 开始写入，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	// 1. 新增管理员
	_, err = session.Insert(user)
	if err != nil {
		session.Rollback()
		return 0, errors.NewSys(err)
	}
	uid = user.Uid // 刚刚生成的管理员ID

	// 2. 新增管理员、用户组关联
	roleIdArr := strings.Split(roleIds, ",") // 逗号分隔
	for _, v := range roleIdArr {
		roleId, err := strconv.Atoi(v)
		if err != nil || roleId <= 0 {
			session.Rollback()
			return 0, errors.NewNormal("参数role_ids格式错误")
		}

		roleUserMD := &RoleUser{
			RoleId: roleId,
			Uid:    uid,
		}

		_, err = session.Insert(roleUserMD)
		if err != nil {
			session.Rollback()
			return 0, errors.NewSys(err)
		}
	}

	err = session.Commit()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return
}

// 更新管理员
func (u *User) Update(user *User, roleIds string) error {
	// 整理数据
	if user.Pwd != "" {
		salt := utils.NewRandomString(5)

		user.Pwd = utils.Md5(utils.Md5(user.Pwd) + salt) // md5两次，第二次带上salt
		user.Salt = salt
	}

	user.UpdateTime = time.Now().Unix()

	// 验证管理员是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(user.Uid).Get(new(User))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("管理员不存在或已删除")
	}

	// 判断管理员用户名是否已存在
	has, err = engine.Where("name=?", user.Name).And("uid!=?", user.Uid).Get(new(User))
	if err != nil {
		return errors.NewSys(err)
	}
	if has {
		return errors.NewNormal("管理员名称已存在")
	}

	// 开始更新，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	// 1. 更新管理员
	_, err = session.ID(user.Uid).Update(user)
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 更新管理员、用户组关联
	// 2.1 删除之前的关联
	_, err = session.Where("uid=?", user.Uid).Delete(new(RoleUser))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2.2 新增关联
	roleIdArr := strings.Split(roleIds, ",") // 逗号分隔
	for _, v := range roleIdArr {
		roleId, err := strconv.Atoi(v)
		if err != nil || roleId <= 0 {
			session.Rollback()
			return errors.NewNormal("参数role_ids格式错误")
		}

		roleUserMD := &RoleUser{
			RoleId: roleId,
			Uid:    user.Uid,
		}

		_, err = session.Insert(roleUserMD)
		if err != nil {
			session.Rollback()
			return errors.NewSys(err)
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
