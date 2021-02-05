package variable

const (
	ParamError          = "Request parameter error"             // 请求参数错误
	FailedSendEmail     = "Failed to send Email"                // 邮件发送失败
	FailedVerifyCode    = "Failed to send verification code"    // 验证码发送失败
	SuccessVerifyCode   = "Send verification code successfully" // 验证码发送成功
	VerifyCodeExpired   = "Verification code has expired"       // 验证码已过期
	VerifyCodeIncorrect = "Incorrect verification code"         // 验证码不正确
	FailedGetUserName   = "Failed to get user name"             // 获取用户姓名失败
	FailedCreateUser    = "Failed to create user"               // 用户添加失败
	SuccessCreateUser   = "User create successfully"            // 用户成功添加
	FailedDeleteUser    = "Failed to delete user"               // 用户删除失败
	SuccessDeleteUser   = "User deleted successfully"           // 用户成功删除
)
