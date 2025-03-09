package helpers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func GenerateFileAvatarName(userId int64, filelName string) string {
	ext := filepath.Ext(filelName)
	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()

	return fmt.Sprintf("%d-%s-%s%s", userId, timestamp, uniqueID, ext)
}



func GetUserKeyProfile(key string) string {
	return fmt.Sprintf("user:%s:profile", key)
}


func GetUserKeyRole(key string) string {
	return fmt.Sprintf("user:%s:role", key)
}