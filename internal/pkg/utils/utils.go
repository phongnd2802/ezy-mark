package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GenerateCliTokenUUID(userId int64) string {
	newUUID := uuid.New()
	// convert UUID to string, remove -
	uuidString := strings.ReplaceAll((newUUID).String(), "", "")
	// 10clitokenijkasdmfasikdjfpomgasdfgl,masdl;gmsdfpgk
	return strconv.FormatInt(userId, 10) + "clitoken" + uuidString
}
