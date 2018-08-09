package controller

import (
	"admin/app/models/backstage"
	"github.com/gin-gonic/gin"
	"strconv"
	"unicode/utf8"
)

// 权限节点
type NodeController struct {
	BaseController
}

func (n *NodeController) Router(e *gin.Engine) {
	group := e.Group("/node")
	{
		group.GET("/list", n.List)
		group.GET("/get", n.Get)
		group.POST("/add", n.Add)
		group.POST("/update", n.Update)
	}
}

// 节点列表
func (n *NodeController) List(ctx *gin.Context) {
	// 调用model
	modelList, _, err := new(backstage.Node).ListAll(nil)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	n.Put(ctx, "list", modelList)

	// 返回
	n.RespOK(ctx)
	return
}

// 获取节点
func (n *NodeController) Get(ctx *gin.Context) {
	// 获取参数
	id, err := n.GetInt(ctx, "id")
	if err != nil || id < 0 {
		n.RespErr(ctx, "参数id格式错误")
		return
	}

	// 调用model
	node, err := new(backstage.Node).Get(id)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	n.Put(ctx, "node", node)

	// 返回
	n.RespOK(ctx)
	return
}

// 新增节点
func (n *NodeController) Add(ctx *gin.Context) {
	// 获取参数
	pid, err := n.GetInt(ctx, "pid")
	if err != nil || pid < 0 {
		n.RespErr(ctx, "参数pid格式错误")
		return
	}

	title := n.GetString(ctx, "title")
	if strLen := utf8.RuneCountInString(title); strLen < 1 || strLen > 10 {
		n.RespErr(ctx, "参数name格式错误")
		return
	}

	nType, err := n.GetInt(ctx, "type")
	if err != nil || !(nType == 1 || nType == 2) {
		n.RespErr(ctx, "参数type格式错误")
		return
	}

	var menuUrl, menuIcon string
	var menuType int
	if nType == 1 { //菜单相关
		menuUrl = n.GetString(ctx, "menu_url")
		menuIcon = n.GetString(ctx, "menu_icon")

		menuType, err = n.GetInt(ctx, "menu_type")
		if err != nil || !(menuType == 1 || menuType == 2) {
			n.RespErr(ctx, "参数menu_type格式错误")
			return
		}
	}

	belongSuper, err := n.GetInt(ctx, "belong_super")
	if err != nil || !(belongSuper == 0 || belongSuper == 1) {
		n.RespErr(ctx, "参数belong_super格式错误")
		return
	}

	weight, err := n.GetInt(ctx, "weight")
	if err != nil || weight < 0 {
		n.RespErr(ctx, "参数weight格式错误")
		return
	}

	// 调用model
	node := &backstage.Node{
		Pid:         pid,
		Title:       title,
		Weight:      weight,
		Type:        nType,
		MenuUrl:     menuUrl,
		MenuIcon:    menuIcon,
		MenuType:    menuType,
		BelongSuper: belongSuper,
	}

	id, err := new(backstage.Node).Add(node)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	n.Put(ctx, "id", id)

	// 返回
	n.RespOK(ctx)
	return
}

//更新
func (n *NodeController) Update(ctx *gin.Context) {
	params := make(map[string]interface{})

	// 获取参数
	id, err := n.GetInt(ctx, "id")
	if err != nil || id < 0 {
		n.RespErr(ctx, "参数id格式错误")
		return
	}

	if v, ok := n.GetParam(ctx, "title"); ok {
		params["title"] = v
	}

	if v, ok := n.GetParam(ctx, "type"); ok {
		if nType, err := strconv.Atoi(v); err != nil || !(nType == 1 || nType == 2) {
			n.RespErr(ctx, "参数type格式错误")
			return
		}

		params["type"] = v
	}

	if v, ok := n.GetParam(ctx, "menu_url"); ok {
		params["menu_url"] = v
	}

	if v, ok := n.GetParam(ctx, "menu_icon"); ok {
		params["menu_icon"] = v
	}

	if v, ok := n.GetParam(ctx, "menu_type"); ok {
		if menuType, err := strconv.Atoi(v); err != nil || !(menuType == 1 || menuType == 2) {
			n.RespErr(ctx, "参数menu_type格式错误")
			return
		}

		params["menu_type"] = v
	}

	if v, ok := n.GetParam(ctx, "belong_super"); ok {
		if belongSuper, err := strconv.Atoi(v); err != nil || !(belongSuper == 0 || belongSuper == 1) {
			n.RespErr(ctx, "参数belong_super格式错误")
			return
		}

		params["belong_super"] = v
	}

	if v, ok := n.GetParam(ctx, "states"); ok {
		if states, err := strconv.Atoi(v); err != nil || !(states == 0 || states == 1) {
			n.RespErr(ctx, "参数states格式错误")
			return
		}

		params["states"] = v
	}

	if v, ok := n.GetParam(ctx, "weight"); ok {
		params["weight"] = v
	}

	// 调用model
	err = new(backstage.Node).Update(id, params)
	if err != nil {
		n.RespErr(ctx, err)
		return
	}

	n.Put(ctx, "id", id)

	// 返回
	n.RespOK(ctx)
	return
}
