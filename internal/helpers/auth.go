package helpers

import (
	"fmt"
	"strings"
)

func ExtractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

func GetUserKeySession(token string) string {
	return fmt.Sprintf("user:%s:session", token)
}

func GetUserKeyOtp(key string) string {
	return fmt.Sprintf("user:%s:otp", key)
}

func GetUserKeyToken(key string) string {
	return fmt.Sprintf("%s:accesstoken", key)
}