package models

import (
	"math"

	"github.com/shopspring/decimal"
)

type BaseModel struct {
}

// 分页列表
type ModelList struct {
	PageIndex int         `json:"page_index"`
	PageSize  int         `json:"page_size"`
	PageCount int         `json:"page_count"`
	Total     int         `json:"total"`
	Items     interface{} `json:"items"`
}

// 计算分页
func (b *BaseModel) Paging(pageIndex, pageSize, total int) (offset int, modelList *ModelList) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (pageIndex - 1) * pageSize

	modelList = &ModelList{
		PageIndex: pageIndex,
		PageSize:  pageSize,
		PageCount: int(math.Ceil(float64(total) / float64(pageSize))),
		Total:     total,
	}

	return
}

func (b *BaseModel) Int64MulInt64By8Bit(ma int64, mb int64) int64 {
	dd := decimal.New(ma, 0)
	dp := decimal.New(mb, 0)

	num := dd.Mul(dp).Div(decimal.New(100000000, 0)).IntPart()
	return num
}

func (b *BaseModel) Int64DivInt64By8Bit(da int64, db int64) int64 {
	dd := decimal.New(da, 0)
	dp := decimal.New(db, 0)

	num := dd.Div(dp).Round(8).Coefficient().Int64()
	return num
}
