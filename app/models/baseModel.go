package models

import "math"

type BaseModel struct {
}

// 计算分页
func (b *BaseModel) Paging(pageIndex, pageSize, total int) (offset, totalPage int) {
	// 计算分页
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	offset = (pageIndex - 1) * pageSize

	totalPage = int(math.Ceil(float64(total / pageSize)))

	return
}
