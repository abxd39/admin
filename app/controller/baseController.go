package controller

type BaseController struct {
	Controller
}

func (b *BaseController) CheckLogin() bool {
	return true
}

func (b *BaseController) GetUid() int {
	return 0
}
