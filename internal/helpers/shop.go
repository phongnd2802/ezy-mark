package helpers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func GenerateShopLogoObjectName(shopName string, filelName string) string {
	ext := filepath.Ext(filelName)
	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()

	return fmt.Sprintf("%s-%s-%s%s", shopName, timestamp, uniqueID, ext)
}


func GenerateBusinessLicenseObjectName(shopName string, filelName string) string {
	ext := filepath.Ext(filelName)
	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()

	return fmt.Sprintf("%s-%s-%s%s", shopName, timestamp, uniqueID, ext)
}