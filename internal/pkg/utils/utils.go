package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/phongnd2802/ezy-mark/internal/global"
)

func GenerateCliTokenUUID(userId int64) string {
	newUUID := uuid.New()
	// convert UUID to string, remove -
	uuidString := strings.ReplaceAll((newUUID).String(), "", "")
	// 10clitokenijkasdmfasikdjfpomgasdfgl,masdl;gmsdfpgk
	return strconv.FormatInt(userId, 10) + "clitoken" + uuidString
}


func GetURLFileName(bucketName string, objectName string) string {
	return fmt.Sprintf("%s/%s/%s", global.Config.MinIOHost, bucketName, objectName)
}

func ExtractFileNameFromURL(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}