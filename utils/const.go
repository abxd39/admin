package utils

const (
	AUTH_NIL    = -1 //取消认证
	AUTH_EMAIL  = 2  //00000010 //邮箱
	AUTH_PHONE  = 1  //00000001	//电话
	AUTH_GOOGLE = 8  //00001000	//google
	AUTH_TWO    = 4  //0100		//二级
	AUTH_FIRST  = 16 //0001 0000 实名认证
)
