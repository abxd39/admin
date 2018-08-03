package models

import (
	"github.com/shopspring/decimal"
)

type BaseModel struct {
	Model
}

type UserInfo struct {
	NickName  string `json:"nick_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
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

func (*BaseModel)Int64MulInt64By8BitString(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	m := dd.Mul(dp)
	d := decimal.New(100000000, 0)
	n := m.Div(d)

	r := n.Div(decimal.New(100000000, 0))
	return r.String()
}

func (*BaseModel)Float64ToInt64By8Bit(s float64) int64 {
	d := decimal.NewFromFloat(s)
	l := d.Round(8).Coefficient().Int64()
	return l
}

func (*BaseModel)Int64ToFloat64By8Bit(b int64) (x float64) {
	a := decimal.New(b, -8)
	x, _ = a.Float64()
	return
}