package models

import (
	"github.com/shopspring/decimal"
	"strconv"
	"fmt"
)

type BaseModel struct {
	Model
}

type UserInfo struct {
	NickName    string  `json:"nick_name"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Status      int     `json:"status"`
	TokenName   string  `json:"token_name"`
	SurplusTrue float64 `xorm:"-" json:"surplus_true"`
	NumTrue     float64 `xorm:"-" json:"num_true"`
}

type SubductionZero struct {
	PriceTrue  float64 `json:"price_true"`
	NumberTrue float64 `json:"number_true"`
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

func (*BaseModel) Int64MulInt64By8BitString(a int64, b int64) string {
	dd := decimal.New(a, 0)
	dp := decimal.New(b, 0)
	m := dd.Mul(dp)
	d := decimal.New(100000000, 0)
	n := m.Div(d)

	r := n.Div(decimal.New(100000000, 0))
	return r.String()
}

func (b *BaseModel) Float64ToInt64By8Bit(s float64) int64 {
	d := decimal.NewFromFloat(s)
	l := d.Round(8).Coefficient().Int64()
	return l
}

func (bb*BaseModel) Int64ToFloat64By8Bit(b int64) (x float64) {
	//fmt.Println("传入的参数值=",b)
	a := decimal.New(b, -8)
	x, _ = a.Float64()
	//fmt.Printf("%v\n",x)
	x = bb.decimal(x)
	//fmt.Println("转换后的结果为=",x)
	return
}

func (b *BaseModel) SubductionZeroMethod(num, price uint64) (rNum, rPrice float64) {
	rNum = b.Int64ToFloat64By8Bit(int64(num))
	rPrice = b.Int64ToFloat64By8Bit(int64(price))
	rPrice = b.decimal(rPrice)
	return
}

func (b *BaseModel) SubductionZeroMethodInt64(num, price int64) (rNum, rPrice float64) {
	rNum = b.Int64ToFloat64By8Bit(num)
	rPrice = b.Int64ToFloat64By8Bit(price)
	rPrice = b.decimal(rPrice)
	return
}

func (b *BaseModel)decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}

