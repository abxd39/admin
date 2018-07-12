package models

import "math"

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
		PageCount: int(math.Ceil(float64(total / pageSize))),
		Total:     total,
	}

	return
}
