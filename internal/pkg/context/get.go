package context

import (
	"context"
	"fmt"
)

func GetSubjectUUID(ctx context.Context) (string, error) {
	sUUID, ok := ctx.Value("subjectUUID").(string)
	if !ok {
		return "", fmt.Errorf("failed to get subject UUID")
	}
	return sUUID, nil
}

