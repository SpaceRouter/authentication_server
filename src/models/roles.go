package models

type Role string

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
