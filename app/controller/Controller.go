package controller

type Controller struct{

}

func (this *Controller)Init() error{
	return nil
}
func (this *Controller) CheckLogin() bool{
	return true
}

func (this *Controller) GetUid() int{

	return 0
}


