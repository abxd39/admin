package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct{}

func (this *CurrencyController) Router(r *gin.Engine) {
	g := r.Group("/currency")
	{
		g.GET("/list", this.GetTradeList)                     //法币挂单列表
		g.GET("/tokens", this.GetTokensList)                  //获取 所有数据货币的名称及货币Id
		g.GET("/order", this.GetOderList)                     //法币成交列表
		g.GET("/total_balance", this.GetTotalCurrencyBalance) //所有法币账户，
		g.GET("/user_detail", this.GetUserDetailList)         //法币账户资产展示
		g.GET("/user_buysell", this.GetBuySellList)           //查看统计买入_卖出_划转
	}
}
func (cu *CurrencyController) GetBuySellList(c *gin.Context) {
	//个人用户交易记录
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Page     int `form:"page" json:"page" binding:"required"`
		Rows     int `form:"rows" json:"rows" `
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	uid := make([]int, 0)
	uid = append(uid, req.Uid)
	fmt.Printf("GetBuySellList%#v\n", uid)
	list, page, total, err := new(models.Order).GetOrderListOfUid(req.Page, req.Rows, req.Uid, req.Token_id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	fmt.Println("000000000000000000000000000000", len(list))
	fmt.Println(list)
	//c.JSON(http.StatusOK, gin.H{"code": 1, "data": list, "msg": err.Error()})
	//return
	type Rsp struct {
		TokenId        uint64
		TokenName      string
		BuyQuantity    float32
		BuyTotalPrice  float32
		SellQuantity   float32
		SellTotalPrice float32
		Transfer       float32
	}
	//法币买入/卖出
	var mapToken map[uint64]*Rsp
	mapToken = make(map[uint64]*Rsp, 0)
	for _, v := range list {

		switch v.AdType {
		case 1: //出售
			fmt.Println("nnnnnnnnnnnnnnnnnnnnnnnnn", v.AdType)
			if _, ok := mapToken[v.TokenId]; ok {
				mapToken[v.TokenId].SellQuantity += float32(v.Num)
				mapToken[v.TokenId].SellTotalPrice += float32(v.Price * v.Num)
			} else {
				mapToken[v.TokenId] = new(Rsp)
				mapToken[v.TokenId].SellTotalPrice = float32(v.Price * v.Num)
				mapToken[v.TokenId].SellQuantity = float32(v.Num)
				mapToken[v.TokenId].TokenId = v.TokenId
			}
		case 2: //购买
			fmt.Println("vvvvvvvvvvvvvvvvvvvvv", v.AdType)
			if _, ok := mapToken[v.TokenId]; ok {
				mapToken[v.TokenId].BuyQuantity += float32(v.Num)
				mapToken[v.TokenId].BuyTotalPrice += float32(v.Price * v.Num)
			} else {
				mapToken[v.TokenId] = new(Rsp)
				mapToken[v.TokenId].BuyQuantity = float32(v.Num)
				mapToken[v.TokenId].BuyTotalPrice = float32(v.Price * v.Num)
				mapToken[v.TokenId].TokenId = v.TokenId
			}
		default:
			continue
		}
	}
	tokenlist, err := new(models.Tokens).GetTokenList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	for _, v := range mapToken {
		for _, t := range tokenlist {
			if v.TokenId == uint64(t.Id) {
				v.TokenName = t.Name
				break
			}
		}

	}
	//
	//if
	returnlist := make([]Rsp, 0)
	for _, value := range mapToken {
		returnlist = append(returnlist, *value)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "total": total, "page": page, "data": returnlist, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetUserDetailList(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Page     int `form:"page" json:"page" binding:"required"`
		Rows     int `form:"rows" json:"rows" `
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}

	result, page, total, err := new(models.UserCurrency).GetCurrencyList(req.Page, req.Rows, req.Uid, req.Token_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": result, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetTotalCurrencyBalance(c *gin.Context) {
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Start_t  string `form:"start_t" json:"start_t" `
		End_t    string `form:"end_t" json:"end_t" `
		Status   int    `form:"status" json:"status" ` //用户账号状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	type listCurrency struct {
		Uid             uint64
		Phone           string
		Email           string
		Status          int     //用户状态
		NickName        string  //用户名
		RegisterTime    int64   //注册时间
		CurrencyBalance float64 //法币账户总和
	}
	list := make([]listCurrency, 0)
	uid := make([]uint64, 0)
	//第一步 调用获取 用户资料
	result, page, total, err := new(models.UserGroup).GetAllUser(req.Page, req.Page_num, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	fmt.Printf("GetTotalCurrencyBalance%#v\n", result)
	for _, v := range result {

		rsp := listCurrency{
			Uid:             v.Uid,
			Phone:           v.Phone,
			Email:           v.Email,
			Status:          v.Status,
			NickName:        v.NickName,
			RegisterTime:    v.RegisterTime,
			CurrencyBalance: 0,
			//法币账户总和为空
		}
		uid = append(uid, v.Uid)
		list = append(list, rsp)
	}
	currlist, err := new(models.UserCurrency).GetAll(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	//找出所有相同的uid 的资产

	for _, v := range uid {
		for _, v1 := range currlist {
			if v == v1.Uid {
				//
			}
		}
	}
	//根据法币账户的成交uid 去获取 用户资料

	// result, total, err := new(models.Ads).xxx(req)
	// if reerr != nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": list, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetTradeList(c *gin.Context) {
	req := models.Currency{}
	fmt.Printf("11111111111111111111111111")
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	fmt.Println(".0.0.0.0.0.0.0.0.0.00.0.0.0.00.0.0....0.0.0.0.0.0")
	fmt.Println(req)
	result, page, total, err := new(models.Ads).GetAdsList(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": result, "msg": "成功"})
	return

}

func (cu *CurrencyController) GetTokensList(c *gin.Context) {
	fmt.Println("tttttttttttttttttttttttttttttttttttttttt")
	list, err := new(models.Tokens).GetTokenList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetOderList(c *gin.Context) {
	//参数一大堆
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Start_t  string `form:"start_t" json:"start_t" `
		End_t    string `form:"end_t" json:"end_t" `
		Status   int    `form:"status" json:"status" `     //订单状态
		Token_id int    `form:"token_id" json:"token_id" ` //货币名称
		Ad_type  int    `form:"adtype" json:"adtype" `     //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	list, page, toal, err := new(models.Order).GetOrderList(req.Page, req.Page_num, req.Ad_type, req.Status, req.Token_id, req.Start_t, req.End_t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": toal, "data": list, "msg": "成功"})
	return
}
