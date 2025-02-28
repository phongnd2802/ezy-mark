package response


var msg = map[int]string{
	ErrCodeSuccess: "success",
	ErrCodeInvalidParams: "invalid params",

	// Auth Message
	ErrCodePendingVerification: "email is being pendding verification",
	ErrCodeEmailAlreadyExists: "email already exists",
	ErrCodeExpiredSession: "expired verification session",
	ErrCodeOtpDoesNotMatch: "otp does not match",
	ErrCodeAuthenticationFailed: "email or password is incorrect!",
	ErrCodeAccountNotVerified: "account is not verified",
	ErrCodeInternalServer: "internal server error",
}

