package consts

const (
	// Authentication
	OTP_EXPIRED_TIME int64 = 60 // 1 minute
	JWT_ACCESS_TOKEN_EXPIRED_TIME string = "1h"
	JWT_REFRESH_TOKEN_EXPIRED_TIME string = "168h" // 7 days

	// Bucket name
	BucketUserAvatar string = "user-avatar"
)
