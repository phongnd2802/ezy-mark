package response

const (
	ErrCodeSuccess = 20000

	ErrCodeInvalidParams = 40000
	ErrCodeBadRequest    = 40001
	// Auth Code
	ErrCodePendingVerification  = 41000
	ErrCodeEmailAlreadyExists   = 41001
	ErrCodeExpiredSession       = 41002
	ErrCodeOtpDoesNotMatch      = 41003
	ErrCodeAuthenticationFailed = 41004
	ErrCodeAccountNotVerified   = 41005
	ErrCodeTokenFlagged         = 41006
	ErrCodeTokenInvalid         = 41007
	ErrCodeUnauthorized         = 41008
	ErrCodeInternalServer       = 50000

	// Shop Code
	ErrCodeShopEmailAlreadyExists = 42000
)
