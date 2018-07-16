package models

import (
	"math"

	"github.com/shopspring/decimal"
)

type BaseModel struct {
}

// 列表
type ModelList struct {
	IsPage    bool        `json:"is_page"`    // 是否分页
	PageIndex int         `json:"page_index"` // 当前页码
	PageSize  int         `json:"page_size"`  // 每页数据条数
	PageCount int         `json:"page_count"` // 总页数
	Total     int         `json:"total"`      // 总数据条数
	Items     interface{} `json:"items"`      // 数据数组
}

// 分页列表
func (b *BaseModel) Paging(pageIndex, pageSize, total int) (offset int, modelList *ModelList) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (pageIndex - 1) * pageSize

	modelList = &ModelList{
		IsPage:    true,
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

// 不分页列表，也包装成一个ModelList
func (b *BaseModel) NoPaging(total int, list interface{}) *ModelList {
	return &ModelList{
		IsPage:    false,
		PageIndex: 1,
		PageSize:  total,
		PageCount: 1,
		Total:     total,
		Items:     list,
	}

}
