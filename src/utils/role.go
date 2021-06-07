package utils

import (
	"github.com/spacerouter/authentication_server/models"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const RolesDir = "/etc/sr/roles"
const DefaultRole = models.Role("NoRole")

func GetRolePath(name string) string {
	return RolesDir + "/" + name + ".yaml"
}

func GetRolePermissions(role models.Role) ([]models.Permission, error) {
	file, err := ioutil.ReadFile(GetRolePath(role.GetName()))
	if err != nil {
		return nil, err
	}

	var roleFile models.RoleFile
	err = yaml.Unmarshal(file, &roleFile)
	if err != nil {
		return nil, err
	}

	return roleFile.Permissions, err
}

func ChangeRoleConfig(role models.Role, permissions []models.Permission) error {

	roleFile := models.RoleFile{
		Permissions: permissions,
	}

	err := WriteRoleConfig(role, roleFile)
	if err != nil {
		return err
	}

	return nil
}

func WriteRoleConfig(role models.Role, roleFile models.RoleFile) error {
	file, err := os.OpenFile(GetRolePath(role.GetName()), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func(open *os.File) {
		_ = open.Close()
	}(file)

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	marshal, err := yaml.Marshal(&roleFile)
	if err != nil {
		return err
	}

	log.Println(string(marshal))
	_, err = file.WriteString(string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func GetRole(roleName string) models.Role {
	return models.Role(roleName)
}

func ListRoles() ([]string, error) {
	files, err := ioutil.ReadDir(RolesDir)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") && !f.IsDir() {
			roles = append(roles, f.Name())
		}
	}
	return roles, nil
}

func IsARole(group string, roleNames []string) bool {
	for _, roleName := range roleNames {
		if group == roleName {
			return true
		}
	}
	return false
}
