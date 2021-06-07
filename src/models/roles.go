package models

type Role string
type Permission string

const (
	ChangeUserInfo Role = "chgusrinfo"
)

func HasRole(roles []Role, role Role) bool {
	for _, uRole := range roles {
		if role == uRole {
			return true
		}
	}
	return false
}

type RoleFile struct {
	Permissions []Permission
}

func PermissionsToStrings(permissions []Permission) []string {
	var result []string
	for _, permission := range permissions {
		result = append(result, string(permission))
	}
	return result
}

func (r *Role) GetName() string {
	return string(*r)
}
