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
		g.GET("/list", this.GetTradeList)
	}
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
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "total": total, "data": result, "msg": "成功"})
	return

}
