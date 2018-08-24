package controller

import (
	"admin/app/models"
	"admin/constant"
	"admin/utils"
	"admin/utils/convert"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

type WallectController struct {
	BaseController
}

func (w *WallectController) Router(r *gin.Engine) {
	g := r.Group("/wallect")
	{
		g.GET("/in_out_trend", w.InOutTrend)
	}
}

// 充币提币走势
func (w *WallectController) InOutTrend(ctx *gin.Context) {
	// 筛选
	filter := make(map[string]interface{})
	if v := w.GetString(ctx, "token_id"); v != "" {
		filter["token_id"] = v
	}
	if v := w.GetString(ctx, "date_begin"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			w.RespErr(ctx, "参数date_begin格式错误")
			return
		}

		filter["date_begin"] = v
	}
	if v := w.GetString(ctx, "date_end"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			w.RespErr(ctx, "参数date_end格式错误")
			return
		}

		filter["date_end"] = v
	}

	// 调用model
	list, err := new(models.TokenInoutDailySheet).InOutTrendList(filter)
	if err != nil {
		w.RespErr(ctx, err)
		return
	}

	// 组装数据
	listLen := len(list)
	x := make([]string, listLen)
	yIn := make([]string, listLen)
	yOut := make([]string, listLen)

	allInTotal := "0"  // 充币总计
	allOutTotal := "0" // 提币总计
	for k, v := range list {
		datetime, _ := time.Parse(utils.LAYOUT_DATE_TIME, v.Date)
		x[k] = datetime.Format("0102")
		yIn[k], _ = convert.StringTo8Bit(v.InTotal)
		yOut[k], _ = convert.StringTo8Bit(v.OutTotal)

		allInTotal, _ = convert.StringAddString(allInTotal, v.InTotal)
		allOutTotal, _ = convert.StringAddString(allOutTotal, v.OutTotal)
	}
	allInTotalFloat, _ := convert.StringTo8Bit(allInTotal)
	allOutTotalFloat, _ := convert.StringTo8Bit(allOutTotal)

	// 设置返回数据
	w.Put(ctx, "x", x)
	w.Put(ctx, "y_in", yIn)
	w.Put(ctx, "y_out", yOut)
	w.Put(ctx, "all_in_total", allInTotalFloat)
	w.Put(ctx, "all_out_total", allOutTotalFloat)

	// 返回
	w.RespOK(ctx)
	return
}
