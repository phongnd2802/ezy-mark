package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// TwMerge combines Tailwind classes and handles conflicts
// Later classes overide earlier ones with the same base
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns a class if a condition is true, otherwise an empty string
// Ex: "bg-green-500", true -> "bg-green-500", false -> ""
func TwIf(class string, condition bool) string {
	result := templ.KV(class, condition)
	if result.Value == true {
		return result.Key
	}
	return ""
}

// mergeAttributes merges multiple Attributes into one
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 16)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	return base64.StdEncoding.EncodeToString(nonceBytes), nil
}

func GenerateCliTokenUUID(userId int64) string {
	newUUID := uuid.New()
	// convert UUID to string, remove -
	uuidString := strings.ReplaceAll((newUUID).String(), "", "")
	// 10clitokenijkasdmfasikdjfpomgasdfgl,masdl;gmsdfpgk
	return strconv.FormatInt(userId, 10) + "clitoken" + uuidString
}
