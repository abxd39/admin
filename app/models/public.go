package models

type Public struct{}

func (p *Public) ParamLegality(page, rows *int) {

	if *page <= 1 {
		*page = 1
	}
	if *rows <= 0 {
		*rows = 50
	}
	return
}

//页码计算
func (p *Public) CalculatePage(total, rows int) (page int) {
	page = total / rows
	if v := total % rows; v != 0 {
		page += 1
	}
	return
}
