package controller

import (
	"unicode/utf8"

	"admin/app/models"
	"admin/app/models/backstage"

	"github.com/gin-gonic/gin"
)

// 用户组
type RoleController struct {
	BaseController
}

func (r *RoleController) Router(e *gin.Engine) {
	group := e.Group("/role")
	{
		group.GET("/list", r.List)
		group.GET("/get", r.Get)
		group.POST("/add", r.Add)
		group.POST("/update", r.Update)
		group.POST("/delete", r.Delete)
	}
}

// 用户组列表
func (r *RoleController) List(ctx *gin.Context) {
	// 获取参数
	isPage, err := r.GetBool(ctx, "is_page", true)
	if err != nil {
		r.RespErr(ctx, "参数is_page格式错误")
		return
	}

	page, err := r.GetInt(ctx, "page", 1)
	if err != nil {
		r.RespErr(ctx, "参数page格式错误")
		return
	}

	rows, err := r.GetInt(ctx, "rows", 10)
	if err != nil {
		r.RespErr(ctx, "参数rows格式错误")
		return
	}

	// 调用model
	var modelList *models.ModelList
	if isPage {
		modelList, _, err = new(backstage.Role).List(page, rows)
	} else {
		modelList, _, err = new(backstage.Role).ListAll()
	}
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	r.Put(ctx, "list", modelList)

	// 返回
	r.RespOK(ctx)
	return
}

// 用户组详情
func (r *RoleController) Get(ctx *gin.Context) {
	// 获取参数
	id, err := r.GetInt(ctx, "id")
	if err != nil || id < 1 {
		r.RespErr(ctx, "参数id格式错误")
		return
	}

	// 调用model
	roleMD := new(backstage.Role)
	role, err := roleMD.Get(id)
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 获取绑定的节点ID
	nodeIds, err := roleMD.GetBindNodeIds(id)
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	r.Put(ctx, "role", role)
	r.Put(ctx, "bind_node_ids", nodeIds)

	// 返回
	r.RespOK(ctx)
	return
}

// 新增用户组
func (r *RoleController) Add(ctx *gin.Context) {
	// 获取参数
	name := r.GetString(ctx, "name")
	if strLen := utf8.RuneCountInString(name); strLen < 1 || strLen > 10 {
		r.RespErr(ctx, "参数name格式错误")
		return
	}

	desc := r.GetString(ctx, "desc", "")

	nodeIds := r.GetString(ctx, "node_ids", "")

	// 调用model
	role := &backstage.Role{
		Name: name,
		Desc: desc,
	}
	id, err := new(backstage.Role).Add(role, nodeIds)
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	r.Put(ctx, "id", id)

	// 返回
	r.RespOK(ctx)
	return
}

// 更新用户组
func (r *RoleController) Update(ctx *gin.Context) {
	params := make(map[string]interface{})

	// 获取参数
	id, err := r.GetInt(ctx, "id")
	if err != nil || id < 1 {
		r.RespErr(ctx, "参数id格式错误")
		return
	}

	if name, ok := r.GetParam(ctx, "name"); ok {
		if strLen := utf8.RuneCountInString(name); strLen < 1 || strLen > 10 {
			r.RespErr(ctx, "参数name格式错误")
			return
		}

		params["name"] = name
	}

	if desc, ok := r.GetParam(ctx, "desc"); ok {
		params["desc"] = desc
	}

	if nodeIds, ok := r.GetParam(ctx, "node_ids"); ok {
		params["node_ids"] = nodeIds
	}

	// 调用model
	err = new(backstage.Role).Update(id, params)
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	r.Put(ctx, "id", id)

	// 返回
	r.RespOK(ctx)
	return
}

// 删除用户组
func (r *RoleController) Delete(ctx *gin.Context) {
	// 获取参数
	id, err := r.GetInt(ctx, "id")
	if err != nil || id < 1 {
		r.RespErr(ctx, "参数id格式错误")
		return
	}

	// 调用model
	err = new(backstage.Role).Delete(id)
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	r.Put(ctx, "id", id)

	// 返回
	r.RespOK(ctx)
	return
}
