package controller

import (
	models "admin/app/models/currency"
	m "admin/app/models/token"
	u "admin/app/models/user"
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
		g.GET("/total_balance", this.GetTotalCurrencyBalance) //所有法币账户，及单个用户的所有交易
		g.GET("user_detail", this.GetUserDetailList)          //用户的法币交易明细
	}
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
	result, total, err := new(u.UserGroup).GetAllUser(req.Page, req.Page_num, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
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
	//根据法币账户的成交uid 去获取 用户资料

	// result, total, err := new(models.Ads).xxx(req)
	// if reerr != nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"code": 0, "total": total, "data": list, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetUserDetailList(c *gin.Context) {
	req := struct {
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	fmt.Println(".0.0.0.0.0.0.0.0.0.00.0.0.0.00.0.0....0.0.0.0.0.0")
	fmt.Println(req)
	// result, total, reerr := new(models.Ads).xxx(req)
	// if reerr != nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": reerr.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"code": 0, "total": total, "data": result, "msg": "成功"})
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
	result, total, reerr := new(models.Ads).GetAdsList(req)
	if reerr != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": reerr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "total": total, "data": result, "msg": "成功"})
	return

}

func (cu *CurrencyController) GetTokensList(c *gin.Context) {
	fmt.Println("tttttttttttttttttttttttttttttttttttttttt")
	list, err := new(m.Tokens).GetTokensList()
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
	list, toal, oerr := new(models.Order).GetOrderList(req.Page, req.Page_num, req.Ad_type, req.Status, req.Token_id, req.Start_t, req.End_t)
	if oerr != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": oerr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "total": toal, "data": list, "msg": "成功"})
	return
}
