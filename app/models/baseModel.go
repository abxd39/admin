package models

import (
	"github.com/shopspring/decimal"
)

type BaseModel struct {
	Model
	UserInfo
}
type UserInfo struct {
	NickName string`json:"nick_name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Status int `json:"status"`
	TokenName string `json:"token_name"`
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
