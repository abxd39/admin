package controller

type BaseController struct {
	Controller
}

func (b *BaseController) GetUid() int {
	return 0
}
