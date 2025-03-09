package helpers

import "fmt"

func HasValidRole(roles []string, userRoles []string) bool {
	for _, role := range roles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
}

func HasValidPermission(permissions []string, userPermissions []string) bool {
	for _, permission := range permissions {
		for _, userPermission := range userPermissions {
			if permission == userPermission {
				return true
			}
		}
	}
	return false
}

func GetKeyRole(key string) string {
	return fmt.Sprintf("role:%s:permissions", key)
}
