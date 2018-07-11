package constant

const (
	// 正常响应
	RESPONSE_CODE_OK = 0

	// 常规错误
	RESPONSE_CODE_ERROR = 90000

	// 系统故障
	RESPONSE_CODE_SYSTEM = 90500

	// 登录会话无效或已掉线
	RESPONSE_CODE_SESSION_INVALID = 90600

	// 登录会话被踢
	RESPONSE_CODE_SESSION_KICK = 90601
)
