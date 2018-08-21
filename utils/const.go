package utils

const (
	AUTH_NIL    = -1 //取消认证
	AUTH_EMAIL  = 2  //00000010 //邮箱
	AUTH_PHONE  = 1  //00000001	//电话
	AUTH_GOOGLE = 8  //00001000	//google
	AUTH_TWO    = 4  //0100		//二级
	AUTH_FIRST  = 16 //0001 0000 实名认证
)

//是否设置资金密码状态标识
const (
	AUTH_TRADEMARK               = 1  //0001资金密码设置状态
	APPLY_FOR_FIRST              = 2  //实名认证申请状态
	APPLY_FOR_SECOND             = 4  //二级认证申请状态
	APPLY_FOR_SECOND_NOT_ALREADY = 8  //二级认证没有通过
	APPLY_FOR_FIRST_NOT_ALREADY  = 16 //一级认证状态未通过
)

const(
	VERIFY_OUT_TOKEN_MARK = 5 // 审核通过
	VERIFY_REVOKE_TOKEN_MARK=6 //审核撤销
)

