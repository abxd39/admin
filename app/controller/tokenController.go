package controller

import (
	"admin/app/models"
	"admin/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
	BaseController
}

func (this *TokenController) Router(r *gin.Engine) {
	g := r.Group("/token")
	{
		g.GET("/list", this.GetTokenOderList)            //bibi p4-1-0 币币委托管理 挂单信息
		g.POST("evacuate_order", this.EvacuateOder)      // p4-1-0 币币委托管理 挂单信息 撤单
		g.GET("/record_list", this.GetRecordList)        //bibi p4-1-1 成交记录
		g.GET("/total_balance", this.GetTokenBalance)    //bibi p2-3-2币币账户统计列表
		g.GET("/user_token_detail", this.GetTokenDetail) //p2-3-2-1查看币币账户资产
		g.GET("/token_cash_list", this.GetTokenCashList) //p4-1-2币兑管理
		g.POST("/delete_cash", this.DeleteCash)          //删除币兑
		g.GET("/modify_cash", this.ModifyCash)           //修改币兑
		g.POST("/add_cash", this.AddCash)                //添加币兑
		g.GET("/change_detail", this.ChangeDetail)       //p2-3-4币币账户变更详情
		g.GET("/fee_list", this.GetFeeInfoList)          //p5-1-0-1币币交易手续费明细
		g.GET("add_take_list", this.GetAddTakeList)      //p5-1-1-1提币手续费明细
	}
}

//p5-1-1-1提币手续费明细
func (this *TokenController) GetAddTakeList(c *gin.Context) {
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Rows     int    `form:"rows" json:"rows" `
		Token_id int    `form:"token_id" json:"token_id" ` //币种
		Uid      int    `form:"uid" json:"uid" `
		Date     uint64 `form:"date",json:"date"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
}

//p5-1-0-1币币交易手续费明细
func (this *TokenController) GetFeeInfoList(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Opt  int    `form:"opt" json:"opt" ` //交易方向
		Uid  int    `form:"uid" json:"uid" `
		Date uint64 `form:"date",json:"date"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.Trade).GetFeeInfoList(req.Page, req.Rows, req.Uid, req.Opt, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
}

func (this *TokenController) DeleteCash(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.ConfigQuenes).DeleteCash(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *TokenController) ModifyCash(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	result, err := new(models.ConfigQuenes).ModifyCash(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", result)

	this.RespOK(c)
	return
}

func (this *TokenController) AddCash(c *gin.Context) {
	req := models.ConfigQuenes{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.ConfigQuenes).AddCash(&req)
	if err != nil {
		fmt.Println("0cccccccccccccccccc0", err)
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

//订单操作 撤单处理

func (this *TokenController) EvacuateOder(c *gin.Context) {
	req := struct {
		Uid      int    `form:"uid" json:"uid" binding:"required"`
		OerderId string `form:"oid" json:"oid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.EntrustDetail).EvacuateOder(req.Uid, req.OerderId)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *TokenController) ChangeDetail(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Date   uint64 `form:"date" json:"date" `
		Search string `form:"search" json:"search" ` //刷选
		Type   int    `form:"type" json:"type" `     //交易方向 买 卖 划转
		Status int    `form:"status" json:"status" ` //用户状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//把货币Id转换为货币名称
	tokenlist, err := new(models.Tokens).GetTokenList()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	//没写完
	//在此分道扬镳
	if req.Search != `` || req.Status != 0 {
		list, err := new(models.UserGroup).GetAllUser(req.Page, req.Rows, req.Status, req.Search)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		value, Ok := list.Items.([]models.UserGroup)
		if !Ok {
			this.RespErr(c, errors.New("assert type UserGroup failed!!"))
			return
		}
		uidlist := make([]uint64, 0)
		for _, v := range value {
			uidlist = append(uidlist, v.Uid)
		}
		monerylist, err := new(models.MoneyRecord).GetMoneyList(req.Page, req.Rows, uidlist)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		tokenValue, ok := monerylist.Items.([]models.MoneyRecord)
		if !ok {
			this.RespErr(c, errors.New("assert type MoneyRecord failed!!"))
			return
		}
		for i, _ := range tokenValue {
			for _, v := range value {
				if int(v.Uid) == tokenValue[i].Uid {
					tokenValue[i].NickName = v.NickName
					tokenValue[i].Email = v.Email
					tokenValue[i].Phone = v.Phone
					tokenValue[i].Status = v.Status
					break
				}
			}
			for _, tv := range tokenlist {
				if tv.Id == tokenValue[i].TokenId {
					tokenValue[i].TokenName = tv.Name
					break
				}
			}
		}
		monerylist.Items = tokenValue
		this.Put(c, "list", monerylist)
		this.RespOK(c)
		return
	} else {
		fmt.Println("1111111111111111111111111", err)
		list, err := new(models.MoneyRecord).GetMoneyListForDateOrType(req.Page, req.Rows, req.Type, req.Date)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		uidList := make([]uint64, 0)
		Value, ok := list.Items.([]models.MoneyRecord)
		if !ok {
			this.RespErr(c, err)
			return
		}
		for _, v := range Value {
			uidList = append(uidList, uint64(v.Uid))
		}
		ulist, err := new(models.UserGroup).GetUserListForUid(uidList)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		for i, _ := range Value {
			for _, v := range ulist {
				if Value[i].Uid == int(v.Uid) {
					Value[i].NickName = v.NickName
					Value[i].Phone = v.Phone
					Value[i].Email = v.Email
					Value[i].Status = v.Status
					break
				}
			}
			for _, vt := range tokenlist {
				if vt.Id == Value[i].TokenId {
					Value[i].TokenName = vt.Name
					break
				}

			}

		}
		list.Items = Value
		this.Put(c, "list", list)
		this.RespOK(c)
		return
	}
}

func (this *TokenController) GetTokenCashList(c *gin.Context) {
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows" `
		Id   int `form:"id" json:"id" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.ConfigQuenes).GetTokenCashList(req.Page, req.Rows, req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c, "成功")
	return
}

func (this *TokenController) GetTokenDetail(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//bibi账户余额
	list, err := new(models.UserToken).GetTokenDetailOfUid(req.Uid, req.Token_id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//bibi 账户统计表
func (this *TokenController) GetTokenBalance(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Status int    `form:"status" json:"status" `
		Search string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	fmt.Printf("GetTokenBalance%#v\n", req)
	list, err := new(models.PersonalProperty).TotalUserBalance(req.Page, req.Rows, req.Status, req.Search)
	if err != nil {
		this.RespErr(c, err)
		return
	}

	this.Put(c, "list", list)
	this.RespOK(c, "成功")
	return
}

//bibi 成交记录
func (this *TokenController) GetRecordList(c *gin.Context) {
	req := struct {
		Page       int    `form:"page" json:"page" binding:"required"`
		Page_num   int    `form:"rows" json:"rows" `
		Uid        int    `form:"uid" json:"uid" `
		Trade_id   int    `form:"trade_id" json:"trade_id" ` //交易类型id 市价交易or 限价交易
		Start_t    string `form:"start_t" json:"start_t" `
		Trade_duad int    `form:"trade_duad" json:"trade_duad" ` //交易对
		Ad_id      int    `form:"ad_id" json:"ad_id" `           //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.Trade).GetTokenRecordList(req.Page, req.Page_num, req.Trade_id, req.Trade_duad, req.Ad_id, req.Uid, req.Start_t)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//币币挂单列表
func (this *TokenController) GetTokenOderList(c *gin.Context) {
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Uid      int    `form:"uid" json:"uid" `
		Trade_id string `form:"trade_id" json:"trade_id" ` //交易类型id 市价交易or 限价交易
		Start_t  int    `form:"start_t" json:"start_t" `
		Symbo    string `form:"symbo" json:"symbo" `  //交易对
		Ad_id    int    `form:"ad_id" json:"ad_id" `  //买卖方向
		Status   int    `form:"status" json:"staus" ` //订单状态
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.EntrustDetail).GetTokenOrderList(req.Page, req.Page_num, req.Ad_id, req.Status, req.Start_t, req.Uid, req.Symbo, req.Trade_id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}
