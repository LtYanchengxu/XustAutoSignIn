package router

import "XustAutoSignIn/variable"

func baseApiRouter() {
	// Ping: GET /api/ping
	variable.Engine.GET("/api/ping", Ping)
	// Get Verify Code: GET /api/verifyCode
	variable.Engine.GET("/api/verifyCode", VerificationCode)
}

func xustAutoSignInApi() {
	xust := variable.Engine.Group("/api/xust")
	// Xust Create User: PUT /api/xust/user
	xust.PUT("user", XUSTCreateUser)
	// Xust Delete User: DELETE /api/xust/user
	xust.DELETE("user", XUSTDeleteUser)
}
