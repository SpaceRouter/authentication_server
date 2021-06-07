package utils

import (
	"fmt"
	"github.com/spacerouter/authentication_server/models"
	"github.com/spacerouter/pam"
	"os/user"
)

// GetUserRole Get user role
func GetUserRole(username string) (*models.Role, error) {
	lookup, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}

	groupIds, err := lookup.GroupIds()
	if err != nil {
		return nil, fmt.Errorf("cannot reach group Ids : %s", err.Error())
	}

	roles, err := ListRoles()
	if err != nil {
		return nil, fmt.Errorf("cannot list roles : %s", err.Error())
	}

	for _, group := range groupIds {
		group, err := user.LookupGroupId(group)
		if err != nil {
			return nil, fmt.Errorf("cannot convert id to name %s : %s", group, err.Error())
		}

		if IsARole(group.Name, roles) {
			role := models.Role(group.Name)
			return &role, nil
		}
	}

	role := DefaultRole
	return &role, nil
}

// GetUserByRole Get information from users belonging to the role
func GetUserByRole(role models.Role) ([]models.UserInfo, error) {
	users, err := pam.ListUsers()
	if err != nil {
		return nil, err
	}

	var usersByRole []models.UserInfo

	for _, userInfo := range users {
		uRole, err := GetUserRole(userInfo.Username)
		if err != nil {
			return nil, err
		}

		if *uRole == role {
			users = append(users, userInfo)
		}
	}
	return usersByRole, nil
}

// UpdateUserPermissions Update Linux permissions to users belonging to the role
func UpdateUserPermissions(role models.Role, permissions []models.Permission) error {
	users, err := GetUserByRole(role)
	if err != nil {
		return err
	}

	groups := models.PermissionsToStrings(permissions)
	groups = append(groups, role.GetName())

	for _, userInfo := range users {
		err := pam.ChangeGroups(userInfo.Login, groups)
		if err != nil {
			return err
		}
	}

	return nil
}

// ChangeUserRole Change user role
func ChangeUserRole(username string, role models.Role) error {
	userRole, err := GetUserRole(username)
	if err != nil {
		return err
	}

	if *userRole != DefaultRole {
		err = pam.RemoveGroup(username, userRole.GetName())
		if err != nil {
			return err
		}
	}

	err = pam.AddGroup(username, role.GetName())
	if err != nil {
		return err
	}

	permissions, err := GetRolePermissions(role)
	if err != nil {
		return err
	}

	err = UpdateUserPermissions(role, permissions)
	if err != nil {
		return err
	}

	return nil
}

func ApplyPolicy() error {
	roles, err := ListRoles()
	if err != nil {
		return err
	}

	for _, roleName := range roles {
		role := GetRole(roleName)

		permissions, err := GetRolePermissions(role)
		if err != nil {
			return err
		}

		err = UpdateUserPermissions(role, permissions)
		if err != nil {
			return err
		}
	}

	return nil
}
