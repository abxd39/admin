package backstage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"admin/app/models"
	"admin/errors"
	"admin/session"
	"admin/utils"

	"github.com/gin-gonic/gin"
)

// 管理员
type User struct {
	models.BaseModel `xorm:"-"`
	Uid              int    `xorm:"uid pk autoincr" json:"uid"`
	Name             string `xorm:"name" json:"name"`
	NickName         string `xorm:"nick_name" json:"nick_name"`
	Pwd              string `xorm:"pwd" json:"-"`
	Salt             string `xorm:"salt" json:"-"`
	States           int    `xorm:"states" json:"states"`
	Remark           string `xorm:"remark" json:"remark"`
	CreateTime       int64  `xorm:"create_time" json:"create_time"`
	UpdateTime       int64  `xorm:"update_time" json:"update_time"`
	LastLoginTime    int64  `xorm:"last_login_time" json:"last_login_time"`
	IsSuper          int    `xorm:"is_super" json:"is_super"`
}

// 管理员 + 管理员所在的用户组名称(多个用逗号分隔)
type UserWithRoleName struct {
	User     `xorm:"extends"`
	RoleName string `json:"role_name"`
}

// 表名
func (*User) TableName() string {
	return "user"
}

// 登录
func (u *User) Login(name, pwd string) (bool, string, int, error) {
	engine := utils.Engine_backstage
	fmt.Println("login")
	user := &User{}
	_, err := engine.Where("name=?", name).Get(user)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		fmt.Println("login", err.Error())
		return false, "", 0, err
	}
	//find 如果不存在 数据库是否会返回 一个错误给我
	if user.States == 0 {
		return false, "", 0, errors.New("该用户已锁定")
	}
	fmt.Printf("数据库中的has值为%s", user.Pwd)
	if utils.Md5(utils.Md5(pwd)+user.Salt) != user.Pwd {
		return false, "", 0, errors.New("密码不对！！")
	}

	// 是否超管
	var isSuper bool
	if user.IsSuper == 1 {
		isSuper = true
	}

	return isSuper, user.NickName, user.Uid, nil
}

// 管理员列表
func (u *User) List(pageIndex, pageSize int, filter map[string]string) (modelList *models.ModelList, list []UserWithRoleName, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("uid")
	query.Alias("u").
		Join("LEFT", []string{new(RoleUser).TableName(), "ru"}, "ru.uid=u.uid").
		Join("LEFT", []string{new(Role).TableName(), "r"}, "r.id=ru.role_id")

	// 筛选
	query.Where("u.is_super=0") // 不显示超管
	if v, ok := filter["phone"]; ok {
		query.And("u.phone like '%?%'", v)
	}

	tempQuery := *query
	count, err := tempQuery.Count(&User{})
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := u.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	err = query.Select("u.*, GROUP_CONCAT(r.name) role_name").Limit(modelList.PageSize, offset).GroupBy("u.uid").Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
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

// 用户组绑定的节点ID
func (u *User) GetBindRoleIds(uid int) (roleIds []int, err error) {
	engine := utils.Engine_backstage
	err = engine.Table(new(RoleUser)).Where("uid=?", uid).Cols("role_id").Find(&roleIds)
	if err != nil {
		return nil, errors.NewSys(err)
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
	user.States = 1
	user.CreateTime = now
	user.UpdateTime = now

	// 判断管理员名称是否已存在
	engine := utils.Engine_backstage
	has, err := engine.Where("name=?", user.Name).Exist(new(User))
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
	if roleIds != "" {
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
	}

	err = session.Commit()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return
}

// 更新管理员
func (u *User) Update(uid int, params map[string]interface{}) error {
	// 验证管理员是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(uid).Exist(new(User))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("管理员不存在或已删除")
	}

	// 整理数据
	userData := make(map[string]interface{})
	if v, ok := params["name"]; ok {
		// 判断管理员用户名是否已存在
		has, err = engine.Where("name=?", v).And("uid!=?", uid).Exist(new(User))
		if err != nil {
			return errors.NewSys(err)
		}
		if has {
			return errors.NewNormal("管理员名称已存在")
		}

		userData["name"] = v
	}

	if v, ok := params["nick_name"]; ok {
		userData["nick_name"] = v
	}

	if v, ok := params["pwd"]; ok {
		salt := utils.NewRandomString(5)

		userData["pwd"] = utils.Md5(utils.Md5(v.(string)) + salt) // md5两次，第二次带上salt
		userData["salt"] = salt
	}

	userData["update_time"] = time.Now().Unix()

	// 开始更新，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	// 1. 更新管理员
	_, err = session.Table(u).ID(uid).Update(userData)
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 更新管理员、用户组关联
	if v, ok := params["role_ids"]; ok {
		roleIds := v.(string)

		// 2.1 删除之前的关联
		_, err = session.Where("uid=?", uid).Delete(new(RoleUser))
		if err != nil {
			session.Rollback()
			return errors.NewSys(err)
		}

		// 2.2 新增关联
		if roleIds != "" {
			roleIdArr := strings.Split(roleIds, ",") // 逗号分隔
			for _, v := range roleIdArr {
				roleId, err := strconv.Atoi(v)
				if err != nil || roleId <= 0 {
					session.Rollback()
					return errors.NewNormal("参数role_ids格式错误")
				}

				roleUserMD := &RoleUser{
					RoleId: roleId,
					Uid:    uid,
				}

				_, err = session.Insert(roleUserMD)
				if err != nil {
					session.Rollback()
					return errors.NewSys(err)
				}
			}
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

// 删除管理员
func (u *User) Delete(uid int) error {
	// 验证用户组是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(uid).Exist(new(User))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("管理员不存在或已删除")
	} else if u.IsSuper == 1 {
		return errors.NewNormal("超管不可删除")
	}

	// 开始删除，事务
	session := engine.NewSession()
	defer session.Close()

	// 1. 删除管理员
	_, err = session.ID(uid).Delete(new(User))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 删除用户组、管理员关联
	_, err = session.Where("uid=?", uid).Delete(new(RoleUser))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	err = session.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

// 检查管理员api权限
func (u *User) CheckPermission(ctx *gin.Context, uid int, api string) (bool, error) {
	// 判断是否超管
	isSuper, err := session.IsSuper(ctx)
	if err != nil {
		return false, err
	}
	if isSuper { // 超管，直接返回true
		return true, nil
	}

	// 判断是否拥有该节点的权限
	// 表名
	userTable := u.TableName()
	roleUserTable := new(RoleUser).TableName()
	roleNodeTable := new(RoleNode).TableName()
	nodeTable := new(Node).TableName()
	nodeAPITable := new(NodeAPI).TableName()

	result := &struct {
		Cnt int
	}{}

	engine := utils.Engine_backstage
	engine.SQL(fmt.Sprintf("SELECT COUNT(n.id) as cnt"+
		" FROM %s u"+
		" JOIN %s ru ON ru.uid=u.uid"+
		" JOIN %s rn ON rn.role_id=ru.role_id"+
		" JOIN %s n ON n.id=rn.node_id"+
		" JOIN %s na ON na.node_id=n.id"+
		" WHERE n.belong_super=0 AND u.uid=%d AND na.api='%s'", userTable, roleUserTable, roleNodeTable, nodeTable, nodeAPITable, uid, api)).Get(result)

	if result.Cnt > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 获取管理员拥有的左侧菜单
func (u *User) MyLeftMenu(ctx *gin.Context) ([]Node, error) {
	// 表名
	userTable := u.TableName()
	roleUserTable := new(RoleUser).TableName()
	roleNodeTable := new(RoleNode).TableName()
	nodeTable := new(Node).TableName()

	var list []Node

	// 判断是否超管
	engine := utils.Engine_backstage
	isSuper, err := session.IsSuper(ctx)
	if err != nil {
		return nil, err
	}
	if isSuper { // 超管，直接返回所有左侧菜单
		engine.Where("type=1").And("states=1").And("menu_type=1").Desc("weight").Find(&list)
	} else {
		uid, err := session.GetUid(ctx)
		if err != nil {
			return nil, err
		}

		engine.SQL(fmt.Sprintf("SELECT n.*"+
			" FROM %s u"+
			" JOIN %s ru ON ru.uid=u.uid"+
			" JOIN %s rn ON rn.role_id=ru.role_id"+
			" JOIN %s n ON n.id=rn.node_id"+
			" WHERE u.uid=%d"+
			" AND n.type=1"+
			" AND n.belong_super=0"+
			" AND n.states=1"+
			" AND n.menu_type=1"+
			" ORDER BY n.weight DESC", userTable, roleUserTable, roleNodeTable, nodeTable, uid)).Find(&list)
	}

	return list, nil
}

// 获取管理员拥有的右侧菜单
func (u *User) MyRightMenu(ctx *gin.Context, pid int) ([]Node, error) {
	// 获取上级信息
	parent, err := new(Node).Get(pid)
	if err != nil {
		return nil, err
	}

	// 表名
	userTable := u.TableName()
	roleUserTable := new(RoleUser).TableName()
	roleNodeTable := new(RoleNode).TableName()
	nodeTable := new(Node).TableName()

	var list []Node

	// 判断是否超管
	engine := utils.Engine_backstage
	isSuper, err := session.IsSuper(ctx)
	if err != nil {
		return nil, err
	}
	if isSuper { // 超管，直接返回所有左侧菜单
		engine.Where("type=1").And("states=1").And("menu_type=2").And(fmt.Sprintf("full_id LIKE '%s%%'", parent.FullId)).Desc("weight").Find(&list)
	} else {
		uid, err := session.GetUid(ctx)
		if err != nil {
			return nil, err
		}

		engine.SQL(fmt.Sprintf("SELECT n.*"+
			" FROM %s u"+
			" JOIN %s ru ON ru.uid=u.uid"+
			" JOIN %s rn ON rn.role_id=ru.role_id"+
			" JOIN %s n ON n.id=rn.node_id"+
			" WHERE u.uid=%d"+
			" AND n.type=1"+
			" AND n.belong_super=0"+
			" AND n.states=1"+
			" AND n.menu_type=2"+
			" AND n.full_id LIKE '%s%%'"+
			" ORDER BY n.weight DESC", userTable, roleUserTable, roleNodeTable, nodeTable, uid, parent.FullId)).Find(&list)
	}

	return list, nil
}
